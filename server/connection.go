package server

import (
	"bytes"
	"log"
	"net"
)

type connection struct {
	tcpConn        *net.TCPConn
	outputBuffer   bytes.Buffer
	internalBuffer []byte
}

func newConnection() (*connection, error) {
	c := connection{}
	c.internalBuffer = make([]byte, 1024)

	addr, err := net.ResolveTCPAddr("tcp4", "localhost:22")
	if err != nil {
		log.Panic(err)
	}

	c.tcpConn, err = net.DialTCP("tcp4", nil, addr)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *connection) buffOutput() error {
	n, err := c.tcpConn.Read(c.internalBuffer)
	c.outputBuffer.Write(c.internalBuffer[:n])
	return err
}

func (c *connection) waitForOutput() error {
	var err error
	for err == nil {
		err = c.buffOutput()
	}
	return err
}
