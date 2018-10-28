package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

var (
	conn   *net.TCPConn
	buffer bytes.Buffer
)

func fetchData() {
	time.Sleep(time.Millisecond * 250)
	for {
		log.Println("Fetching data over HTTP...")

		resp, err := http.Get("http://localhost:8080/")
		if err != nil {
			log.Panicln(err)
		}

		io.Copy(conn, resp.Body)
		time.Sleep(time.Millisecond * 500)
	}
}

func sendData() {
	for {
		log.Println("Sending data over HTTP...")

		_, err := http.Post("http://localhost:8080/", "application/octet-stream", &buffer)
		if err != nil {
			log.Panicln(err)
		}

		time.Sleep(time.Millisecond * 500)
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
	conn = c

	go sendData()
	go fetchData()

	b := make([]byte, 1024)
	for {
		n, _ := conn.Read(b)
		buffer.Write(b[:n])
	}
}
