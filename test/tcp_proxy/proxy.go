package main

import (
	"io"
	"log"
	"net"
)

func handleConn(conn *net.TCPConn) {
	addr, err := net.ResolveTCPAddr("tcp4", "mc.hypixel.net:25565")
	if err != nil {
		log.Panic(err)
	}

	sshConn, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		log.Panic(err)
	}

	go io.Copy(sshConn, conn)
	io.Copy(conn, sshConn)
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

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Panic(err)
		}

		go handleConn(conn)
	}
}
