package data_structure_test

import (
	"bufio"
	"os"

	ds "github.com/frealcone/external-sort/data_structure"
)

type LocalFileReader struct {
	Source *bufio.Reader
}

func (fr LocalFileReader) Read() (string, error) {
	line, _, err := fr.Source.ReadLine()
	return string(line), err
}

func (fr *LocalFileReader) ChangeSource(source *os.File) {
	fr.Source = bufio.NewReader(source)
}

func (fr LocalFileReader) Copy() ds.FileReader {
	return &LocalFileReader{fr.Source}
}
