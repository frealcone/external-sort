package data_structure

import "os"

// FileReader 定义了从文件中读取文本, 且一次读取相当于一条记录的文本的方法
type FileReader interface {
	Read() (string, error)
	ChangeSource(source *os.File)
	Copy() FileReader
}
