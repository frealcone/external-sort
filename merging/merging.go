package merging

import (
	"fmt"
	"io"
	"os"

	ds "github.com/frealcone/external-sort/data_structure"
	"github.com/frealcone/external-sort/preprocessing"
	"golang.org/x/sync/errgroup"
)

// MergeRunLength 将K个游程文件合并为同一, readers为文件读取器
// 在无虚游程时, len(readers) == K; 否则, len(readers) < K
func MergeRunLength(K int, readers []ds.FileReader, cvt ds.Convertable, result *os.File) error {
	// 获取实际游程文件数量
	N := len(readers)
	// 创建PK树
	pkTree := ds.NewEmptyRaceTree(K)
	// 创建游程文件读取通道
	rlDatas := make([]chan string, N)
	for i := 0; i < N; i++ {
		rlDatas[i] = make(chan string, 20)
	}
	// 创建游程文件读取goroutine
	egrp := new(errgroup.Group)
	for i := 0; i < N; i++ {
		idx := i
		egrp.Go(func() error {
			return preprocessing.ReadFromFile(readers[idx], rlDatas[idx])
		})
	}
	// 创建记录通道
	records := make([]chan ds.FileRecord, N)
	for i := 0; i < N; i++ {
		records[i] = make(chan ds.FileRecord, 20)
	}
	// 创建文本转化记录goroutine
	for i := 0; i < N; i++ {
		idx := i
		egrp.Go(func() error {
			return preprocessing.SToFR(rlDatas[idx], cvt, records[idx])
		})
	}
	// 获取第一轮数据
	for i := 0; i < N; i++ {
		if record, ok := <-records[i]; ok {
			pkTree.Put(record, i)
		} else {
			pkTree.Put(nil, i)
		}
	}
	// =========================================================
	// DEBUG
	// 开始输出结果
	for {
		r, next := pkTree.Peak()
		if r == nil {
			break
		}
		fmt.Fprint(result, r.String())
		if record, ok := <-records[next]; ok {
			pkTree.Put(record, next)
		} else {
			pkTree.Put(nil, next)
		}
	}
	// =========================================================

	if err := egrp.Wait(); err != io.EOF {
		return err
	}
	return nil
}

// Merge 将所有游程文件合并为单一目标文件
func Merge(rlPath string, runLengths map[string]int, targetFile string, reader ds.FileReader, cvt ds.Convertable) error {
	n := len(runLengths)

	var k int // 合并路数
	if n <= 3 {
		k = n
	} else {
		k = n%3 + 2
	}

	// 建立最佳合并树
	mergeTree := new(ds.MergeTree)
	for rl, w := range runLengths {
		mergeTree.Insert(ds.MergeTreeNode{
			FileName: rl,
			Weight:   w,
		})
	}

	// 打开目标文件, 有则创建, 无则清空
	tf, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	for k < n {
		mg, err := os.CreateTemp(rlPath, "mg_*.rl") // 合并到临时文件
		if err != nil {
			return err
		}

		n -= (k - 1)

		rlnames, _ := mergeTree.Merge(mg.Name(), k) // 执行合并

		readers := make([]ds.FileReader, 0)
		for _, rlname := range rlnames {
			rl, err := os.Open(rlname)
			if err != nil {
				return err
			}
			defer rl.Close()

			r := reader.Copy()
			r.ChangeSource(rl)

			readers = append(readers, r)
		}

		err = MergeRunLength(k, readers, cvt, mg)
		mg.Close()
		if err != nil {
			return err
		}
	}

	// 最后一轮合并
	rlnames, _ := mergeTree.Merge("", k)
	readers := make([]ds.FileReader, 0)
	for _, rlname := range rlnames {
		rl, err := os.Open(rlname)
		if err != nil {
			return err
		}
		defer rl.Close()

		r := reader.Copy()
		r.ChangeSource(rl)

		readers = append(readers, r)
	}

	err = MergeRunLength(k, readers, cvt, tf)
	if err != nil {
		return err
	}
	return nil
}
