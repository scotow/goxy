package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

var (
//conn *net.TCPConn
//buffer *bytes.Buffer
)

func fetchData() {
	time.Sleep(time.Millisecond * 1000)
	for {
		resp, err := http.Get("http://localhost:8080/")
		if err != nil {
			log.Panicln(err)
		}

		//io.Copy(conn, resp.Body)
		resp.Body.Close()
		time.Sleep(time.Millisecond * 2000)
	}
}

func sendData(buffer *bytes.Buffer) {
	for {
		/*fmt.Println(buffer.Len())
		_, err := http.Post("http://localhost:8080/", "application/octet-stream", &buffer)
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println(buffer.Len())

		fmt.Println("Sending data over HTTP")*/

		time.Sleep(time.Millisecond * 2000)
		log.Println("Before copy", buffer.Len())
		io.Copy(ioutil.Discard, buffer)
		log.Println("After copy", buffer.Len())
	}
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":2222")
	if err != nil {
		log.Panic(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Panic(err)
	}

	c, err := listener.AcceptTCP()
	if err != nil {
		log.Panic(err)
	}

	var buffer bytes.Buffer

	//conn = c

	//go fetchData()
	go sendData(&buffer)

	io.Copy(&buffer, c)
}
