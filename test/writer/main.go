package main

import (
	"io"
	"os"
)

type Person struct {
	data string
}

func (pe *Person) Read(p []byte) (n int, err error) {
	 return copy(p, pe.data), nil
}

func main() {
	p := &Person{"Hello"}
	io.Copy(os.Stdout, p)
}