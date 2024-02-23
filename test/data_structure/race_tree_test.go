package data_structure_test

import (
	"testing"

	ds "github.com/frealcone/external-sort/data_structure"
)

func TestNewRaceTree(t *testing.T) {
	mockData := []ds.FileRecord{Integer(23), Integer(17), Integer(7), Integer(60)}
	ans := []ds.FileRecord{Integer(7), Integer(17), Integer(7)}

	tree := ds.NewRaceTree(mockData)

	for i := 0; i < len(ans); i++ {
		if tree[i].CompareTo(ans[i]) != 0 {
			t.Errorf("race doesn't match with answer at %d, race tree: %v\n", i, tree)
		}
	}
}

func TestPut(t *testing.T) {
	mockData := []ds.FileRecord{Integer(23), Integer(17), Integer(7), Integer(60)}
	tree := ds.NewRaceTree(mockData)

	if r, _ := tree.Peak(); r.CompareTo(Integer(7)) != 0 {
		t.Fatalf("new race tree failed\n")
	}

	tree.Put(Integer(1), 3)

	if r, _ := tree.Peak(); r.CompareTo(Integer(1)) != 0 {
		t.Fatalf("Put() failed\n")
	}
}
