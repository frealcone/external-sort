package external_sort

import (
	"os"

	ds "github.com/frealcone/external-sort/data_structure"
	"github.com/frealcone/external-sort/merging"
	"github.com/frealcone/external-sort/preprocessing"
)

// ExtSort 方法将source文件中的内容排序, 并最终输出到destination文件中, N代表内存最大可以存放N条记录
// 注: 参数r用于从文件中读取文本以用于转化为记录, 一次读取操作仅读出一条记录
// 另外, 用户实现的Read方法需要自己决定缓冲区大小, ExtSort传入该方法的切片大小为0
func ExtSort(r ds.FileReader, cvt ds.Convertable, destination string, N int) error {
	// 预处理
	flengths, path, err := preprocessing.Preprocess(r, cvt, N)
	if err != nil {
		return err
	}

	// 合并预处理结果
	err = merging.Merge(path, flengths, destination, r, cvt)
	if err != nil {
		return err
	}

	// 清理中间文件
	os.RemoveAll(path)

	return nil
}
