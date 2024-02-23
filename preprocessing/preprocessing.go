package preprocessing

import (
	"io"
	"os"

	ds "github.com/frealcone/external-sort/data_structure"
	"golang.org/x/sync/errgroup"
)

// 文件记录的预处理阶段, 默认将游程文件置于系统设置的临时文件存放目录
func Preprocess(reader ds.FileReader, cvt ds.Convertable, N int) (map[string]int, string, error) {
	// 游程文件所在的目录
	path, err := os.MkdirTemp("/tmp", "rl_*")
	if err != nil {
		return nil, "", err
	}

	// 创建文本管道
	messages := make(chan string, 20)
	// 创建记录管道
	records := make(chan ds.FileRecord, 20)

	// 创建管理goroutine报错的errgroup
	egrp := new(errgroup.Group)

	// 启动文件读取goroutine
	egrp.Go(func() error {
		return ReadFromFile(reader, messages)
	})

	// 启动文本转化记录goroutine
	egrp.Go(func() error {
		return SToFR(messages, cvt, records)
	})

	// 启动游程写入goroutine
	var r map[string]int
	egrp.Go(func() error {
		var err error
		r, err = InsertRecord(records, path, N)
		return err
	})

	err = egrp.Wait()
	if err != nil {
		return nil, "", err
	}

	return r, path, nil
}

// 从文件中读取文本, 并写入到output
func ReadFromFile(r ds.FileReader, output chan string) error {
	msg, err := r.Read()
	for err == nil {
		output <- msg
		msg, err = r.Read()
	}

	close(output)

	if err == io.EOF {
		return nil
	}
	return err
}

// 将文本转化为记录
func SToFR(input chan string, cvt ds.Convertable, output chan ds.FileRecord) error {
	defer close(output)

	for msg := range input {
		r, err := cvt.Convert(msg)
		if err != nil {
			return err
		}

		output <- r
	}

	return nil
}

// InsertRecord 将input中的数据有序地输出到游程文件中
func InsertRecord(input chan ds.FileRecord, outputPath string, N int) (map[string]int, error) {
	memory := ds.NewFRecordSpace(
		ds.WithMaxRecordNum(N),
		ds.WithTempPath(outputPath),
	)

	for record := range input {
		err := memory.Insert(record)
		if err != nil {
			return nil, err
		}
	}

	return memory.Flush(), nil
}
