package server

import (
	"bytes"
	"errors"
	"net"
	"sync"
)

var (
	ErrInvalidRemoteAddress = errors.New("invalid remote TCP address")
	ErrCannotReachRemote    = errors.New("cannot open TCP connection to remote host")
)

type connection struct {
	tcpConn *net.TCPConn

	lock           sync.RWMutex
	outputBuffer   bytes.Buffer
	internalBuffer []byte

	socketClosed chan bool
	closing      bool
}

func newConnection(address string) (*connection, error) {
	c := connection{}
	c.internalBuffer = make([]byte, 1024)
	c.socketClosed = make(chan bool)

	addr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, ErrInvalidRemoteAddress
	}

	c.tcpConn, err = net.DialTCP("tcp4", nil, addr)
	if err != nil {
		return nil, ErrCannotReachRemote
	}

	return &c, nil
}

func (c *connection) pipeSocketBuffer(channel chan<- error) {
	for {
		// Wait for some data to fill the TCP socket internal buffer.
		n, err := c.tcpConn.Read(c.internalBuffer)

		// Stop on error.
		if err != nil {
			c.tcpConn.Close()
			channel <- err
			return
		}

		// Copy data while locking the buffer.
		c.lock.Lock()
		_, err = c.outputBuffer.Write(c.internalBuffer[:n])
		c.lock.Unlock()

		// Stop on error.
		if err != nil {
			c.tcpConn.Close()
			channel <- err
			return
		}
	}
}

func (c *connection) start() {
	pipeWorker := make(chan error)
	go c.pipeSocketBuffer(pipeWorker)

	select {
	case <-c.socketClosed:
		c.tcpConn.Close()
	case <-pipeWorker:
		c.lock.Lock()
		c.closing = true
		c.lock.Unlock()
	}
}
