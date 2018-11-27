package client

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	. "github.com/scotow/goxy/common"
)

const (
	maximumSizeWriteGet = 128
	defaultUserAgent    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"
)

func Dial(remoteAddr *net.TCPAddr) (*Conn, error) {
	httpAddr := fmt.Sprintf("http://%s/", remoteAddr.String())

	req, err := http.NewRequest("GET", httpAddr, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	token, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	id, err := NewIdFromToken(string(token))
	if err != nil {
		return nil, err
	}

	var readToken, writeToken string
	fmt.Sscanf(resp.Header.Get("X-Referer"), "%s %s", &readToken, &writeToken)

	conn := Conn{id, remoteAddr, readToken, writeToken}

	return &conn, nil
}

type Conn struct {
	id         *Id
	remoteAddr *net.TCPAddr

	readToken  string
	writeToken string
}

func (c *Conn) buildHttpUrl() string {
	return fmt.Sprintf("http://%s/%s", c.remoteAddr.String(), c.id.RandomPath())
}

func (c *Conn) buildHttpUrlWithHider(hider *Hider) string {
	return fmt.Sprintf("http://%s/%s.%s", c.remoteAddr.String(), c.id.RandomPath(), hider.Extension)
}

func (c *Conn) Read(b []byte) (n int, err error) {
	hider, err := RandomHider()
	if err != nil {
		n = 0
		return
	}

	req, err := http.NewRequest("GET", c.buildHttpUrlWithHider(hider), nil)
	if err != nil {
		n = 0
		return
	}

	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Referer", c.readToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		n = 0
		return
	}
	defer resp.Body.Close()

	c.readToken = resp.Header.Get("X-Referer")

	reader := hider.GetExtractor(resp.Body)

	for {
		read, er := reader.Read(b[n:])
		n += read
		err = er

		if err == io.EOF {
			err = nil
			break
		}

		if err != nil {
			break
		}
	}

	//fmt.Fprintf(c.logger, "Read: buffer size: %d. Read: %d.\n", len(b), n)

	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	hider, err := RandomHider()
	if err != nil {
		n = 0
		return
	}

	req := new(http.Request)

	if len(b) < maximumSizeWriteGet {
		req, err = http.NewRequest("GET", c.buildHttpUrlWithHider(hider), nil)
		if err != nil {
			n = 0
			return
		}
		req.Header.Set("Authorization", base64.StdEncoding.EncodeToString(b))
	} else {
		req, err = http.NewRequest("POST", c.buildHttpUrlWithHider(hider), bytes.NewReader(b))
		if err != nil {
			n = 0
			return
		}
	}

	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Referer", c.writeToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		n = 0
		return
	}
	defer resp.Body.Close()

	c.writeToken = resp.Header.Get("X-Referer")

	n = len(b)

	//fmt.Fprintf(c.logger, "Write: buffer size: %d. Written: %d.\n", len(b), n)
	return
}

func (c *Conn) Close() error {
	return nil
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
