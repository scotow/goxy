package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var (
	ErrSessionCreation = errors.New("cannot create session with the server")
)

type connection struct {
	tcpConn  *net.TCPConn
	httpAddr string

	id             string
	outputBuffer   bytes.Buffer
	internalBuffer []byte
}

func newConnection(tcpConn *net.TCPConn, httpAddr string) (*connection, error) {
	c := connection{}
	c.tcpConn = tcpConn
	c.httpAddr = httpAddr
	c.internalBuffer = make([]byte, 1024)

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

func (c *connection) fetchData(interval time.Duration) error {
	for {
		//log.Println("Fetching data over HTTP...")

		resp, err := http.Get("http://localhost:8080/")
		if err != nil {
			return err
		}

		io.Copy(c.tcpConn, resp.Body)
		time.Sleep(interval)
	}
}

func (c *connection) sendData(interval time.Duration) error {
	for {
		//log.Println("Sending data over HTTP...")

		_, err := http.Post("http://localhost:8080/", "application/octet-stream", &c.outputBuffer)
		if err != nil {
			return err
		}

		time.Sleep(interval)
	}
}
