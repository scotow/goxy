package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
)

func main() {
	var b bytes.Buffer

	log.Println(b.Len())
	b.Write(make([]byte, 1000))
	log.Println(b.Len())
	b.ReadByte()
	log.Println(b.Len())
	io.CopyN(ioutil.Discard, &b, 499)
	log.Println(b.Len())
	b.Write(make([]byte, 1000))
	log.Println(b.Len())
	log.Println(b.Bytes())
}
