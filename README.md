# external-sort

`external-sort`是使用`Go.lang`编写的外排序库; 支持对用户自定义文件的内容进行排序.

`external-sort` is an external sort algorithm lib built with `Go.lang` which supports sorting user-defined files.

## 使用 (Get Started)

`external-sort`算法的定义如下:

```go
func ExtSort(r DiskReader, cvt FileRecordConv, destination string, N int) error
```

- 描述

`ExtSort`共接收4个参数: r用户读取文件, cvt用于将文件的字符(或二进制)内容转化为文件记录类型的数据, destination为排好序的数据最终存储到的磁盘文件路径, N为对可用内存所能容纳下的文件记录数量的估计.

其中, 文件读取, 转化相关的接口定义如下, 用户**必须**自行实现这些接口:

1. DiskReader

```go
Read() (string, error)
ChangeSource(source *os.File)
Copy() DiskReader
```

其中, 文件读取接口声明了三个方法: `Read()`方法从待排序文件中, 读取可以转化为一个文件记录类型变量的字符(或二进制)内容, 并将文件指针移动到下一次读取的位置; `ChangeSource()`方法将待排序文件修改为`source`; `Copy()`方法将当前`DiskReader`类型的变量复制一份并返回;

2. FileRecordConv

```go
Convert(s string) (FileRecord, error)
```

`Convert()`方法负责将存放了字符(或二进制)数据的`string`类型字符串转化为文件记录(`FileRecord`)类型的变量;

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

`CompareTo()`将当前`FileRecord`类型变量与另一变量`fr`进行比较, 如当前变量大于/等于/小于`fr`时, 返回值分别小于/等于/大于0; `String()`将当前变量再次变回字符(或二进制)数据;

> **注意**:
> `destination`如果表示的是一个已有的文件, 那么`ExtSort()`会将该文件清空.
> `Convert()`和`String()`互为逆操作, 也就是说, `string`类型的变量`s`在不进行其他操作的情况下, 经过`fr, _ := Convert(s); t := fr.String()`, 得到变量`t`, 则s等于t

### 案例 (Case)

已有文件numbers.txt:

```
10
300
34
89
...
```

下述代码可将numbers.txt进行排序

```go
import (
	"bufio"
	"os"
	"strconv"

	esort "github.com/frealcone/external-sort"
)

type Integer int

func (i Integer) CompareTo(s esort.Sortable) int {
	return int(i) - int(s.(Integer))
}

func (i Integer) String() string {
	return strconv.Itoa(int(i)) + "\n"
}

type IntegerConv struct{}

func (ic IntegerConv) Convert(s string) (esort.FileRecord, error) {
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

func (fr LocalFileReader) Copy() esort.DiskReader {
	return &LocalFileReader{fr.Source}
}

func main() {
	file, _ := os.Open("./mock_data")
	source := bufio.NewReader(file)

	reader := &LocalFileReader{Source: source}

	convertor := new(IntegerConv)

	esort.ExtSort(reader, convertor, "./result", 100)
}
```