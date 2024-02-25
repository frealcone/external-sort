package test

import (
	"bufio"
	"os"
	"testing"

	external_sort "github.com/frealcone/external-sort"
	data_structure_test "github.com/frealcone/external-sort/test/data_structure"
)

func TestExtSort(t *testing.T) {
	sf, err := os.Create("./mock_data")
	if err != nil {
		t.Fatalf("failed to create mock data file\n")
	}

	data_structure_test.GenData(sf, 100000)
	sf.Close()

	sf, err = os.Open("./mock_data")
	if err != nil {
		t.Fatalf("failed to open mock data file")
	}

	reader := new(data_structure_test.LocalFileReader)
	reader.Source = bufio.NewReader(sf)

	cvt := new(data_structure_test.IntegerConv)

	err = external_sort.ExtSort(reader, cvt, "./mock_result", 10000)
	if err != nil {
		t.Fatalf("failed to sort data: %s", err.Error())
	}
}
