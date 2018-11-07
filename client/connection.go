package client

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/scotow/goxy/common"
	log "github.com/sirupsen/logrus"
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
	lock           sync.Mutex
	internalBuffer []byte
	dynamicSleep   *dynamicSleep
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

func (c *connection) askForConnection() error {
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

	log.WithFields(log.Fields{
		"id":         c.id,
		"address":    c.httpAddr,
		"remoteAddr": c.remoteAddr,
	}).Info("Connection created.")

	return nil
}

func (c *connection) start() {
	socketClosed := make(chan error)
	go c.pipeSocketBuffer(socketClosed)

	stopFetch := make(chan bool)
	fetchStopped := make(chan bool)
	go c.fetchData(stopFetch, fetchStopped)

	stopSend := make(chan bool)
	sendStopped := make(chan bool)
	go c.sendData(stopSend, sendStopped)

	select {
	case <-fetchStopped:
		stopSend <- true
		c.tcpConn.Close()
		<-socketClosed
	case <-sendStopped:
		stopFetch <- true
		c.tcpConn.Close()
		<-socketClosed
	case <-socketClosed:
		stopFetch <- true
		stopSend <- true
		http.Get(fmt.Sprintf("http://%s/%s/close", c.httpAddr, c.id))
	}

	log.WithFields(log.Fields{
		"id":         c.id,
		"address":    c.httpAddr,
		"remoteAddr": c.remoteAddr,
	}).Info("Connection closed.")
}

func (c *connection) pipeSocketBuffer(socketClosed chan<- error) {
	for {
		// Wait for some data to fill the TCP socket internal buffer.
		n, err := c.tcpConn.Read(c.internalBuffer)

		// Stop on error.
		if err != nil {
			socketClosed <- err
			return
		}

		// Copy data while locking the buffer.
		c.lock.Lock()
		_, err = c.outputBuffer.Write(c.internalBuffer[:n])
		c.lock.Unlock()

		// Stop on error.
		if err != nil {
			socketClosed <- err
			return
		}
	}
}

func (c *connection) fetchData(shouldStop <-chan bool, receiveStop chan<- bool) {
	for {
		select {
		case <-c.dynamicSleep.sleep():
			// Fetch pending data from remote socket.
			resp, err := http.Get(fmt.Sprintf("http://%s/%s", c.httpAddr, c.id))
			if err != nil {
				receiveStop <- true
				return
			}

			n, err := io.Copy(c.tcpConn, resp.Body)
			if err != nil {
				receiveStop <- true
				return
			}

			if resp.StatusCode == http.StatusGone {
				receiveStop <- true
				return
			}

			// If remote output buffer had nothing, increase next fetch interval.
			if n == 0 {
				c.dynamicSleep.increment()
			} else {
				c.dynamicSleep.reset()
			}
		case <-shouldStop:
			return
		}
	}
}

func (c *connection) sendData(shouldStop <-chan bool, receiveStop chan<- bool) {
	for {
		select {
		case <-c.dynamicSleep.sleepOriginal():
			c.lock.Lock()

			// If the output buffer is empty, don't send it.
			if c.outputBuffer.Len() == 0 {
				c.lock.Unlock()
			} else {
				// Otherwise send it and reset fetch dynamic interval.
				resp, err := http.Post(fmt.Sprintf("http://%s/%s", c.httpAddr, c.id), "application/octet-stream", &c.outputBuffer)
				c.lock.Unlock()

				if err != nil {
					receiveStop <- true
					return
				}

				if resp.StatusCode == http.StatusGone {
					receiveStop <- true
					return
				}
				c.dynamicSleep.reset()
			}
		case <-shouldStop:
			return
		}
	}
}
