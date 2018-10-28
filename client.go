package goxy

import (
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	httpAddr string
}

func NewClient(localPort int, httpAddr, remoteAddr string) (*Client, error) {
	client := Client{httpAddr}
	return &client, nil
}

func (c *Client) CheckServerHealth() bool {
	resp, err := http.Get(c.httpAddr)
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

func WaitUntilServerUp(retryInterval time.Duration) {

}
