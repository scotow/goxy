package main

import (
	"log"
	"net"
	"time"
)

func handleConn(conn *net.TCPConn) {
	conn.Write([]byte("Hello"))
	time.Sleep(time.Second)
	conn.Close()
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":6666")
	if err != nil {
		log.Panic(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Panic(err)
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Panic(err)
		}

		go handleConn(conn)
	}
}
