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

	notifyRead   chan bool
	shouldNotify bool

	clientClosed chan bool
	closing      bool
}

func newConnection(address string) (*connection, error) {
	c := connection{}
	c.internalBuffer = make([]byte, 1024)

	c.notifyRead = make(chan bool)
	c.clientClosed = make(chan bool)

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

func (c *connection) start() {
	socketClosed := make(chan error)
	go c.pipeSocketBuffer(socketClosed)

	select {
	case <-c.clientClosed:
		c.tcpConn.Close()
	case <-socketClosed:
		c.lock.Lock()
		c.closing = true
		c.lock.Unlock()
	}
}

func (c *connection) pipeSocketBuffer(socketClosed chan<- error) {
	for {
		// Wait for some data to fill the TCP socket internal buffer.
		n, err := c.tcpConn.Read(c.internalBuffer)

		// Stop on error.
		if err != nil {
			c.tcpConn.Close()
			socketClosed <- err
			return
		}

		// Copy data while locking the buffer.
		c.lock.Lock()
		shouldNotify := c.shouldNotify
		_, err = c.outputBuffer.Write(c.internalBuffer[:n])
		c.lock.Unlock()

		if shouldNotify {
			c.notifyRead <- true
		}

		// Stop on error.
		if err != nil {
			c.tcpConn.Close()
			socketClosed <- err
			return
		}
	}
}
