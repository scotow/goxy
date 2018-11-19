package server

import (
	"github.com/scotow/goxy/common"
	"io"
	"net"
	"strconv"
	"time"
)

func newConn(localAddr, remoteAddr *net.TCPAddr) (*Conn, error) {
	id := strconv.FormatInt(time.Now().UnixNano(), 10)
	readC, readNC, readEC := make(chan []byte), make(chan int), make(chan error)
	writeC, writeNC, writeEC := make(chan []byte), make(chan int), make(chan error)

	return &Conn{
		id,
		localAddr, remoteAddr, &common.State{},
		readC, readNC, readEC,
		writeC, writeNC, writeEC,
	}, nil
}

type Conn struct {
	id         string
	localAddr  *net.TCPAddr
	remoteAddr *net.TCPAddr
	state      *common.State

	readC  chan []byte
	readNC chan int
	readEC chan error

	writeC  chan []byte
	writeNC chan int
	writeEC chan error
}

func (c *Conn) Read(b []byte) (n int, err error) {
	c.readC <- b

	n = <-c.readNC
	err = <-c.readEC

	//fmt.Fprintf(c.logger, "Read: buffer size: %d. Read: %d.\n", len(b), n)
	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	if c.state.IsClosed() {
		n, err = 0, io.ErrClosedPipe
		return
	}

	for {
		c.writeC <- b[n:]
		n += <-c.writeNC
		err = <-c.writeEC

		if err != nil {
			break
		}

		if n == len(b) {
			break
		}
	}

	//fmt.Fprintf(c.logger, "Write: buffer size: %d. Written: %d.\n", len(b), n)
	return
}

func (c *Conn) Close() error {
	return nil
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
