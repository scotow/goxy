package client2

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Conn struct {
	id string

	remoteAddr *net.TCPAddr
}

func (c *Conn) Read(b []byte) (n int, err error) {
	httpAddr := fmt.Sprintf("http://%s/read/%s", c.remoteAddr.String(), c.id)

	resp, err := http.Post(httpAddr, "*/*", strings.NewReader(strconv.Itoa(len(b))))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	n, _ = resp.Body.Read(b)
	err = nil

	// TODO: Check for error on read (should be EOF).
	//fmt.Printf("Read: buffer size: %d. Read: %d.\n", len(b), n)
	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	httpAddr := fmt.Sprintf("http://%s/write/%s", c.remoteAddr.String(), c.id)

	resp, err := http.Post(httpAddr, "*/*", bytes.NewReader(b))
	if err != nil {
		n = 0
		return
	}
	defer resp.Body.Close()

	n = len(b)
	// TODO: Check for end of file with custom HTTP status code.

	//fmt.Printf("Write: buffer size: %d. Written: %d.\n", len(b), n)
	return
}

func (c *Conn) Close() error {
	panic("implement me")
}

func (c *Conn) LocalAddr() net.Addr {
	panic("implement me")
}

func (c *Conn) RemoteAddr() net.Addr {
	panic("implement me")
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

func Dial(remoteAddr *net.TCPAddr) (*Conn, error) {
	httpAddr := fmt.Sprintf("http://%s/create", remoteAddr.String())

	resp, err := http.Get(httpAddr)
	if err != nil {
		return nil, err
	}

	id, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	conn := Conn{string(id), remoteAddr}
	return &conn, nil
}
