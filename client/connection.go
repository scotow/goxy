package client

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/scotow/goxy/common"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	fetchInterval = time.Millisecond * 50
)

var (
	ErrSessionCreation = errors.New("cannot create session with the server")
)

type connection struct {
	tcpConn    *net.TCPConn
	httpAddr   string
	remoteAddr string
	id         string

	outputBuffer   bytes.Buffer
	bufferLock     sync.Mutex
	internalBuffer []byte
	dynamicSleep   *dynamicSleep
	closed         bool
}

func newConnection(tcpConn *net.TCPConn, httpAddr string, remoteAddr string) *connection {
	c := connection{}
	c.tcpConn = tcpConn
	c.httpAddr = httpAddr
	c.remoteAddr = remoteAddr
	c.internalBuffer = make([]byte, 1024)
	c.dynamicSleep = newDynamicSleep(fetchInterval, 10)

	return &c
}

func (c *connection) AskForConnection() error {
	reqBody := fmt.Sprintf("%s %s %s", common.AppName, common.Version, c.remoteAddr)
	resp, err := http.Post(fmt.Sprintf("http://%s/create", c.httpAddr), "text/plain", strings.NewReader(reqBody))
	if err != nil {
		return ErrSessionCreation
	}

	if resp.StatusCode != http.StatusOK {
		return ErrSessionCreation
	}

	idBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ErrSessionCreation
	}

	c.id = string(idBytes)
	return nil
}

func (c *connection) buffOutput() error {
	// Buff local output while locking access to the buffer.
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
		if c.closed {
			return c.tcpConn.Close()
		}
		err := c.buffOutput()
		if err != nil {
			fmt.Println("sending close")
			c.closed = true
			http.Get(fmt.Sprintf("http://%s/%s/close", c.httpAddr, c.id))
			return err
		}
	}
}

func (c *connection) fetchData() error {
	for {
		if c.closed {
			return nil
		}

		// Fetch pending data from remote socket.
		resp, err := http.Get(fmt.Sprintf("http://%s/%s", c.httpAddr, c.id))
		if err != nil {
			c.closed = true
			return err
		}

		n, err := io.Copy(c.tcpConn, resp.Body)

		if err != nil {
			c.closed = true
			return err
		}

		if resp.StatusCode == http.StatusGone {
			c.closed = true
			return nil
		}

		// If remote output buffer had nothing, increase next fetch interval.
		if n == 0 {
			c.dynamicSleep.sleepIncrement()
		} else {
			c.dynamicSleep.sleepReset()
		}
	}
}

func (c *connection) sendData() error {
	for {
		if c.closed {
			return nil
		}

		c.bufferLock.Lock()

		// If the output buffer is empty, don't send it.
		if c.outputBuffer.Len() == 0 {
			c.bufferLock.Unlock()
			//c.dynamicSleep.sleepOriginal()
		} else {
			// Otherwise send it and reset fetch dynamic interval.
			resp, err := http.Post(fmt.Sprintf("http://%s/%s", c.httpAddr, c.id), "application/octet-stream", &c.outputBuffer)
			c.bufferLock.Unlock()

			if err != nil {
				c.closed = true
				return err
			}

			if resp.StatusCode == http.StatusGone {
				c.closed = true
				return nil
			}
			c.dynamicSleep.sleepReset()
		}
	}
}

func (c *connection) start() {
	go c.waitForOutput()
	go c.sendData()
	c.fetchData()
}
