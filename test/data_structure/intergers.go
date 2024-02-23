package data_structure_test

import (
	"strconv"

	ds "github.com/frealcone/external-sort/data_structure"
)

type Integer int

func (i Integer) CompareTo(s ds.Sortable) int {
	return int(i) - int(s.(Integer))
}

func (i Integer) String() string {
	return strconv.Itoa(int(i)) + "\n"
}

type IntegerConv struct{}

func (ic IntegerConv) Convert(s string) (ds.FileRecord, error) {
	i, e := strconv.Atoi(s)
	if e != nil {
		return nil, e
	}
	return Integer(i), e
}

func IsIntegersSorted(is []Integer) bool {
	for i := 1; i < len(is); i++ {
		if is[i].CompareTo(is[i-1]) < 0 {
			return false
		}
	}
	return true
}
