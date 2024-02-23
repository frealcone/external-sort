package preprocessing_test

import (
	"bufio"
	"os"
	"testing"

	"github.com/frealcone/external-sort/preprocessing"
	data_structure_test "github.com/frealcone/external-sort/test/data_structure"
)

func TestPreprocess(t *testing.T) {
	sf, err := os.Create("./mock_data")
	if err != nil {
		t.Fatalf("failed to create mock data file\n")
	}

	data_structure_test.GenData(sf, 100)
	sf.Close()

	sf, err = os.Open("./mock_data")
	if err != nil {
		t.Fatalf("failed to open mock data file")
	}

	reader := new(data_structure_test.LocalFileReader)
	reader.Source = bufio.NewReader(sf)

	cvt := new(data_structure_test.IntegerConv)

	_, _, err = preprocessing.Preprocess(reader, cvt, 100)
	if err != nil {
		t.Fatalf("failed to sort data: %s\n", err.Error())
	}

	data_structure_test.DeleData(sf)
}
