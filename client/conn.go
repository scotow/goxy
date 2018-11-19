package client

import (
	"bytes"
	"fmt"
	"github.com/scotow/goxy/common"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	stateOpen = iota
	stateClosing
	stateClosed
)

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

	conn := Conn{string(id), remoteAddr, &common.State{}, new(http.Client)}

	return &conn, nil
}

type Conn struct {
	id         string
	remoteAddr *net.TCPAddr
	state      *common.State
	httpClient *http.Client
}

func (c *Conn) Read(b []byte) (n int, err error) {
	if c.state.IsClosed() {
		n, err = 0, io.ErrClosedPipe
		return
	}

	httpAddr := fmt.Sprintf("http://%s/read/%s", c.remoteAddr.String(), c.id)

	resp, err := http.Post(httpAddr, "*/*", strings.NewReader(strconv.Itoa(len(b))))
	if err != nil {
		fmt.Println("HTTP read POST request", err.Error())
		return
	}
	defer resp.Body.Close()

	for {
		read, er := resp.Body.Read(b[n:])
		n += read
		err = er

		if err == io.EOF {
			break
		}

		if err != nil {
			return
		}
	}

	if c.state.IsClosed() {
		err = io.EOF
	} else {
		err = nil
	}

	//fmt.Fprintf(c.logger, "Read: buffer size: %d. Read: %d.\n", len(b), n)

	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	if c.state.IsClosed() {
		n, err = 0, io.ErrClosedPipe
		return
	}

	httpAddr := fmt.Sprintf("http://%s/write/%s", c.remoteAddr.String(), c.id)

	resp, err := http.Post(httpAddr, "*/*", bytes.NewReader(b))
	if err != nil {
		n = 0
		return
	}
	defer resp.Body.Close()

	n = len(b)
	// TODO: Check for end of file with custom HTTP status code.

	//fmt.Fprintf(c.logger, "Write: buffer size: %d. Written: %d.\n", len(b), n)
	return
}

func (c *Conn) Close() error {
	return nil
}

func (c *Conn) WaitForRemoteClose() {
	for {
		if c.state.IsClosed() {
			return
		}

		httpAddr := fmt.Sprintf("http://%s/wait/%s", c.remoteAddr.String(), c.id)

		resp, err := http.Get(httpAddr)
		if err != nil {
			return
		}

		if resp.StatusCode == 200 {
			c.state.SetClosed()
			return
		} else {
			log.Println("Wait for remote close request timed out.", resp.StatusCode)
		}
	}
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
