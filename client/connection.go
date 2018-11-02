package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
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
	tcpConn  *net.TCPConn
	httpAddr string
	id       string

	outputBuffer   bytes.Buffer
	bufferLock     sync.Mutex
	internalBuffer []byte
	dynamicSleep   *dynamicSleep
}

func newConnection(tcpConn *net.TCPConn, httpAddr string) (*connection, error) {
	c := connection{}
	c.tcpConn = tcpConn
	c.httpAddr = httpAddr
	c.internalBuffer = make([]byte, 1024)
	c.dynamicSleep = newDynamicSleep(fetchInterval, 10)

	return &c, nil
}

func (c *connection) AskForConnection() error {
	resp, err := http.Get(fmt.Sprintf("http://%s/create", c.httpAddr))
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

func (c *connection) fetchData() error {
	for {
		resp, err := http.Get(fmt.Sprintf("http://%s/%s", c.httpAddr, c.id))
		if err != nil {
			return err
		}

		n, err := io.Copy(c.tcpConn, resp.Body)
		if n == 0 {
			c.dynamicSleep.sleepIncrement()
		} else {
			c.dynamicSleep.sleepReset()
		}
	}
}

func (c *connection) sendData() error {
	for {
		c.bufferLock.Lock()

		if c.outputBuffer.Len() == 0 {
			c.bufferLock.Unlock()
			c.dynamicSleep.sleepOriginal()
		} else {
			_, err := http.Post(fmt.Sprintf("http://%s/%s", c.httpAddr, c.id), "application/octet-stream", &c.outputBuffer)
			c.bufferLock.Unlock()

			if err != nil {
				return err
			}

			c.dynamicSleep.sleepReset()
		}
	}
}
