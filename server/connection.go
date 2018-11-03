package server

import (
	"bytes"
	"log"
	"net"
	"sync"
)

type connection struct {
	tcpConn *net.TCPConn

	bufferLock     sync.Mutex
	outputBuffer   bytes.Buffer
	internalBuffer []byte
}

func newConnection(address string) (*connection, error) {
	c := connection{}
	c.internalBuffer = make([]byte, 1024)

	addr, err := net.ResolveTCPAddr("tcp4", address)
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
	c.bufferLock.Lock()
	c.outputBuffer.Write(c.internalBuffer[:n])
	c.bufferLock.Unlock()
	return err
}

func (c *connection) waitForOutput() error {
	var err error
	for err == nil {
		err = c.buffOutput()
	}
	return err
}
