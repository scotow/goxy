package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type Client struct {
	localPort  int
	httpAddr   string
	remoteAddr string
}

func NewClient(localPort int, httpAddr, remoteAddr string) (*Client, error) {
	c := Client{}
	c.localPort = localPort
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
	log.Println("Remote server up and running.")
}

func (c *Client) Start() error {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", c.localPort))
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	for {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}

		conn, err := newConnection(tcpConn, c.httpAddr)
		if err != nil {
			continue
		}

		if err := conn.AskForConnection(); err != nil {
			continue
		}

		go func() {
			go conn.waitForOutput()

			interval := time.Millisecond * 500
			go conn.sendData(interval)
			time.Sleep(interval / 2)
			conn.fetchData(interval)
		}()
	}
}
