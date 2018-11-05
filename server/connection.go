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

	bufferLock     sync.Mutex
	outputBuffer   bytes.Buffer
	internalBuffer []byte

	socketClosed  chan bool
	stopBuffering chan bool
	closed        bool
}

func newConnection(address string) (*connection, error) {
	c := connection{}
	c.internalBuffer = make([]byte, 1024)
	c.socketClosed = make(chan bool)
	c.stopBuffering = make(chan bool)

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

func (c *connection) buffOutput() error {
	n, err := c.tcpConn.Read(c.internalBuffer)
	if err != nil {
		return err
	}

	c.bufferLock.Lock()
	defer c.bufferLock.Unlock()

	_, err = c.outputBuffer.Write(c.internalBuffer[:n])
	return err
}

func (c *connection) waitForOutput() error {
	for {
		select {
		case <-c.stopBuffering:
			c.tcpConn.Close()
			return nil
		default:
			err := c.buffOutput()
			if err != nil {
				c.socketClosed <- true
				return err
			}
		}
	}
}
