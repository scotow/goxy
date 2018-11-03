package client

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	localAddr  string
	httpAddr   string
	remoteAddr string
}

func NewClient(localAddr string, httpAddr, remoteAddr string) (*Client, error) {
	c := Client{}
	c.localAddr = localAddr
	c.httpAddr = httpAddr
	c.remoteAddr = remoteAddr

	return &c, nil
}

func (c *Client) CheckServerStatus() bool {
	resp, err := http.Get(fmt.Sprintf("http://%s/%s", c.httpAddr, "status"))
	if err != nil {
		return false
	}

	if resp.StatusCode != http.StatusOK {
		return false
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	bodyString := string(bodyBytes)
	return bodyString == "OK"
}

func (c *Client) WaitUntilServerUp(retryInterval time.Duration) {
	for !c.CheckServerStatus() {
		time.Sleep(retryInterval)
	}
	log.WithField("address", c.httpAddr).Info("Remote server up and running.")
}

func (c *Client) Start() error {
	addr, err := net.ResolveTCPAddr("tcp", c.localAddr)
	if err != nil {
		log.WithFields(log.Fields{
			"address": c.localAddr,
			"error":   err,
		}).Error("Invalid listening address.")
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.WithFields(log.Fields{
			"address": c.localAddr,
			"error":   err,
		}).Error("Cannot start TCP listener.")
		return err
	}

	for {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			log.WithField("address", c.localAddr).Warn("Cannot accept TCP connection.")
			continue
		}

		conn := newConnection(tcpConn, c.httpAddr, c.remoteAddr)

		if err := conn.AskForConnection(); err != nil {
			conn.tcpConn.Close()
			continue
		}

		go conn.start()
	}
}
