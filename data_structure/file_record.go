package data_structure

// Sortable 定义了CompareTo方法, 该方法将本Sortable数据与s比较
// 返回值>0或=0或<0分别代表本数据大于/等于/小于s
type Sortable interface {
	CompareTo(s Sortable) int
}

// FileRecord 是从文件中读取到的文本转化为的数据记录, 一条记录必须满足:
// 1. 可排序(即, 可以与另一条记录比较大小)
// 2. 可转化为文本, 从而输出到结果文件中
type FileRecord interface {
	Sortable
	String() string
}

// Convertable 是将文本转化为一条数据记录的工具
type Convertable interface {
	Convert(s string) (FileRecord, error)
}
