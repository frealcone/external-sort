package data_structure_test

import (
	"fmt"
	"math/rand"
	"os"
)

func GenData(output *os.File, num int) {
	for i := 0; i < num; i++ {
		data := rand.Intn(num)
		fmt.Fprintln(output, data)
	}
}

func DeleData(f *os.File) {
	os.Remove(f.Name())
	f.Close()
}
