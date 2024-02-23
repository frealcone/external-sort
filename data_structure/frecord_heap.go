package data_structure

import (
	"container/heap"
	"fmt"
	"os"
)

// 文件类型的小根堆
type FileRecordHeap []FileRecord

func (frh FileRecordHeap) Len() int {
	return len(frh)
}

func (frh FileRecordHeap) Less(i, j int) bool {
	return frh[i].CompareTo(frh[j]) < 0
}

func (frh FileRecordHeap) Swap(i, j int) {
	frh[i], frh[j] = frh[j], frh[i]
}

func (frh *FileRecordHeap) Pop() interface{} {
	hp := *frh
	n := len(hp)
	r := hp[n-1]
	*frh = hp[:n-1]
	return r
}

func (frh *FileRecordHeap) Push(x interface{}) {
	*frh = append(*frh, x.(FileRecord))
}

// 选择置换排序所用空间
type FRecordSpace struct {
	init          FileRecordHeap // 初始堆
	rl            *os.File       // 游程文件
	temp          FileRecordHeap // 新堆
	N             int            // 最大记录数量
	path          string         // 存放游程文件的路径
	status        int            // 空间状态
	weightCounter map[string]int // 游程记录数
}

func NewFRecordSpace(options ...FRecordSpaceOption) *FRecordSpace {
	r := new(FRecordSpace)

	for _, f := range options {
		f(r)
	}

	// 两个堆永远初始化为nil
	r.init = nil
	r.temp = nil

	r.rl = nil

	r.weightCounter = make(map[string]int)

	return r
}

const (
	statusEmpty    = 0 // 初始状态, 第一批数据还未完全加载进空间
	statusNotFixed = 1 // 新堆尚未调整
	statusFixed    = 2 // 新堆已被调整
)

// 将记录有序插入到游程文件中
func (s *FRecordSpace) Insert(record FileRecord) error {
	var err error = nil

	// 根据空间状态决定记录存放位置
	switch s.status {
	case statusEmpty:
		// 直接将新记录放入初始堆中
		s.init.Push(record)
		// 初始堆满, 调整空间状态
		if s.init.Len() >= s.N {
			s.status = statusNotFixed
			s.rl, err = os.CreateTemp(s.path, "rl_*") // 建立游程
			heap.Init(&s.init)                        // 调整初始堆
		}
	case statusNotFixed:
		or := heap.Pop(&s.init).(FileRecord)
		if or.CompareTo(record) <= 0 { // 比较刚输出的记录与当前记录
			heap.Push(&s.init, record) // 如当前记录大于旧记录, 将其放入初始堆
		} else {
			s.temp.Push(record) // 否则放入新堆
		}
		fmt.Fprint(s.rl, or.String())   // 输出旧记录
		if s.temp.Len() >= (s.N >> 1) { // 如果新堆占用了一半以上的空间
			heap.Init(&s.temp)     // 则调整新堆
			s.status = statusFixed // 调整空间状态
		}
	case statusFixed:
		or := heap.Pop(&s.init).(FileRecord)
		if or.CompareTo(record) <= 0 { // 比较刚输出的记录与当前记录
			heap.Push(&s.init, record) // 如当前记录大于旧记录, 将其放入初始堆
		} else {
			heap.Push(&s.temp, record) // 否则放入新堆
		}
		fmt.Fprint(s.rl, or.String()) // 输出旧记录
		if s.temp.Len() >= s.N {      // 如果新堆占据了全部的空间, 新堆成为初始堆
			s.status = statusNotFixed
			s.init = s.temp
			s.temp = nil
			s.rl.Close()
			s.rl, err = os.CreateTemp(s.path, "rl_*")
		}
	}

	// 记录输出数量
	if (s.status == statusFixed || s.status == statusNotFixed) && err == nil { // 确保本轮操作进行过输出
		if _, ok := s.weightCounter[s.rl.Name()]; !ok {
			s.weightCounter[s.rl.Name()] = 1
		} else {
			s.weightCounter[s.rl.Name()] += 1
		}
	}

	return err
}

// 将记录立刻全部保存到磁盘上
func (s *FRecordSpace) Flush() map[string]int {
	// 记录初始堆输出数量
	if _, ok := s.weightCounter[s.rl.Name()]; !ok {
		s.weightCounter[s.rl.Name()] = s.init.Len()
	} else {
		s.weightCounter[s.rl.Name()] += s.init.Len()
	}
	for s.init.Len() > 0 { // 先保存初始堆数据
		record := heap.Pop(&s.init).(FileRecord)
		fmt.Fprint(s.rl, record.String())
	}
	s.rl.Close()
	rl, err := os.CreateTemp(s.path, "rl_*")
	if err != nil {
		return nil
	}
	defer rl.Close()
	if s.temp.Len() > 0 { // 记录新堆输出数量
		heap.Init(&s.temp)
		s.weightCounter[rl.Name()] = s.temp.Len()
	}
	for s.temp.Len() > 0 { // 再保存新堆
		record := heap.Pop(&s.temp).(FileRecord)
		fmt.Fprint(rl, record.String())
	}

	return s.weightCounter
}

type FRecordSpaceOption func(*FRecordSpace)

func WithMaxRecordNum(N int) FRecordSpaceOption {
	return func(frs *FRecordSpace) {
		frs.N = N
	}
}

func WithTempPath(path string) FRecordSpaceOption {
	return func(frs *FRecordSpace) {
		frs.path = path
	}
}
