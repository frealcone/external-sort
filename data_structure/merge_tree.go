package data_structure

// MergeTreeNode 合并树节点
type MergeTreeNode struct {
	FileName string // 待合并文件名称
	Weight   int    // 待合并文件在合并树中的权重
}

// MergeTree 合并树数据结构 (Huffman Tree)
type MergeTree []MergeTreeNode

// Insert 将合并树节点插入合并树
func (mt *MergeTree) Insert(node MergeTreeNode) {
	huffman := *mt
	for i := 0; i < len(huffman); i++ {
		if huffman[i].Weight >= node.Weight {
			*mt = append(huffman[:i], append([]MergeTreeNode{node}, huffman[i:]...)...)
			return
		}
	}
	*mt = append(huffman, node)
}

// Merge 将合并树中权重最小的K个节点合并, 合并后新的节点名称为nodeName
// 返回被合并的节点名称和实际合并的节点数量
func (mt *MergeTree) Merge(nodeName string, k int) ([]string, int) {
	n := min(len(*mt), k)
	names := make([]string, n)

	w := 0
	for i := 0; i < n; i++ {
		names[i] = (*mt)[i].FileName
		w += (*mt)[i].Weight
	}
	*mt = (*mt)[n:]
	mt.Insert(MergeTreeNode{
		FileName: nodeName,
		Weight:   w,
	})

	return names, n
}
