package data_structure

type RaceTree []FileRecord

// 初始化PK树, 用于初始化的记录数组必须长度为偶数
func NewRaceTree(records []FileRecord) RaceTree {
	n := len(records)

	var tree RaceTree = make([]FileRecord, 2*n-1)
	for i := 0; i < len(tree); i++ {
		tree[i] = nil
	}

	i, j := n-1, 2*n-2
	for i >= 0 {
		tree[j] = records[i]
		tree.Adjust(j)

		i--
		j--
	}

	return tree
}

func NewEmptyRaceTree(N int) RaceTree {
	var tree RaceTree = make([]FileRecord, 2*N-1)
	for i := 0; i < len(tree); i++ {
		tree[i] = nil
	}
	return tree
}

func (tree RaceTree) Adjust(i int) {
	for i > 0 {
		if i&1 != 0 { // 奇数下标
			if tree[i+1] == nil || (tree[i] != nil && tree[i].CompareTo(tree[i+1]) < 0) {
				tree[i/2] = tree[i]
				i >>= 1
			} else {
				tree[i/2] = tree[i+1]
				i >>= 1
			}
		} else { // 偶数下标
			if tree[i-1] == nil || (tree[i] != nil && tree[i].CompareTo(tree[i-1]) < 0) {
				tree[(i-1)/2] = tree[i]
				i = (i - 1) >> 1
			} else {
				tree[(i-1)/2] = tree[i-1]
				i = (i - 1) >> 1
			}
		}
	}
}

func (tree RaceTree) Peak() (FileRecord, int) {
	r := len(tree) / 2
	for i := len(tree) / 2; i < len(tree); i++ {
		if tree[i] == tree[0] {
			r = i
			break
		}
	}
	return tree[0], r - (len(tree) / 2)
}

func (tree RaceTree) Put(record FileRecord, idx int) {
	n := len(tree)
	idx = idx + n/2

	tree[idx] = record
	tree.Adjust(idx)
}
