package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"os"
	"time"
)

var (
	buffer bytes.Buffer
)

func sendData() {
	for {
		time.Sleep(time.Millisecond * 50)
		io.Copy(os.Stdout, &buffer)

		/*log.Printf("Before copy length: %d\n", buffer.Len())
		copied, err := io.Copy(os.Stdout, &buffer)
		log.Println("Copied: ", copied, err)
		log.Printf("After copy length: %d\n", buffer.Len())*/
	}
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":5555")
	if err != nil {
		log.Panic(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Panic(err)
	}

	conn, err := listener.AcceptTCP()
	if err != nil {
		log.Panic(err)
	}

	go sendData()

	b := make([]byte, 16)
	for {
		n, _ := conn.Read(b)
		buffer.Write(b[:n])
	}
}
