package server2

import (
	"net"
	"strconv"
	"time"
)

func newConn(localAddr, remoteAddr *net.TCPAddr) *Conn {
	id := strconv.FormatInt(time.Now().UnixNano(), 10)
	readC, readNC := make(chan []byte), make(chan int)
	writeC, writeNC := make(chan []byte), make(chan int)

	return &Conn{id, localAddr, remoteAddr, readC, readNC, writeC, writeNC}
}

type Conn struct {
	id         string
	localAddr  *net.TCPAddr
	remoteAddr *net.TCPAddr

	readC  chan []byte
	readNC chan int

	writeC  chan []byte
	writeNC chan int
}

func (c *Conn) Read(b []byte) (n int, err error) {
	c.readC <- b
	n = <-c.readNC

	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	c.writeC <- b
	n = <-c.writeNC

	return
}

func (c *Conn) Close() error {
	panic("implement me")
}

func (c *Conn) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *Conn) SetDeadline(t time.Time) error {
	panic("implement me")
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	panic("implement me")
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	panic("implement me")
}