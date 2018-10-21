package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:6000")
	if err != nil {
		log.Panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Panic(err)
	}

	fmt.Fprintln(conn, "Hello, World")
}
