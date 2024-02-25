# external-sort

`external-sort` is an external sort algorithm lib built with `Go.lang` which supports sorting user-defined files.

## Get Started

Defination of `external-sort`:

```go
func ExtSort(r ds.FileReader, cvt ds.Convertable, destination string, N int) error
```

## Description

`ExtSort()` need four args: `r` defines how to read data from the original file; `cvt` converts charaters(or binary data) read from the origin file to variable of `FileRecord` type; `destination` defines the path of the final result; `N` is the maximum number of `FileRecord` variables that memory can contain.

The related interfaces are as follows, user **must** implement those interfaces:

1. FileReader

```go
Read() (string, error)
ChangeSource(source *os.File)
Copy() DiskReader
```

`FileReader` declared four methods: `Read()` reads data that could transform to exactly one `FileRecord` from the original file; `ChangeSource()` changes origin file to `source`; `Copy()` copies current `FileReader` to generate a new `FileReader` and returns it.

2. Convertable

```go
Convert(s string) (FileRecord, error)
```

`Convert()` converts s(who stores charater(or binary data) read from the original file) to `FileRecord`

3. FileRecord

```go
type Sortable interface {
	CompareTo(s Sortable) int
}
```

```go
Sortable
String() string
```

`CompareTo()` compares current current file record with another record `fr`; The `CompareTo` function must return an integer that's less than, equal to, or greater than zero if the current record is considered to be respectively less than, equal to, or greater than `fr`.

> **NOCICE**:
> If `destination` represents an existing file, `ExtSort()` will truncate it, otherwise, `ExtSort` will create it.
> If `fr, _ := Convert(s); t := fr.String()`, `s == t`.

## Case

numbers.txt:

```
10
300
34
89
...
```

How to sort numbers.txt with inadequate memory:

```go
import (
	"bufio"
	"os"
	"strconv"

	esort "github.com/frealcone/external-sort"
	eds "github.com/frealcone/external-sort/data_structure"
)

type Integer int

func (i Integer) CompareTo(s eds.Sortable) int {
	return int(i) - int(s.(Integer))
}

func (i Integer) String() string {
	return strconv.Itoa(int(i)) + "\n"
}

type IntegerConv struct{}

func (ic IntegerConv) Convert(s string) (eds.FileRecord, error) {
	i, e := strconv.Atoi(s)
	if e != nil {
		return nil, e
	}
	return Integer(i), e
}

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

func (fr LocalFileReader) Copy() eds.FileReader {
	return &LocalFileReader{fr.Source}
}

func main() {
	file, _ := os.Open("./numbers.txt")
	source := bufio.NewReader(file)

	reader := &LocalFileReader{Source: source}

	convertor := new(IntegerConv)

	esort.ExtSort(reader, convertor, "./result", 100)
}
```