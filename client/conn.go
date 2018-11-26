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

	resp, err := http.Get(httpAddr)
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

	conn := Conn{id, remoteAddr}

	return &conn, nil
}

type Conn struct {
	id         *Id
	remoteAddr *net.TCPAddr
}

func (c *Conn) buildHttpUrl() string {
	return fmt.Sprintf("http://%s/%s", c.remoteAddr.String(), c.id.RandomPath())
}

func (c *Conn) Read(b []byte) (n int, err error) {
	resp, err := http.Get(c.buildHttpUrl())
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
	req := new(http.Request)

	if len(b) < maximumSizeWriteGet {
		req, err = http.NewRequest("GET", c.buildHttpUrl(), nil)
		if err != nil {
			n = 0
			return
		}
		req.Header.Set("Authorization", base64.StdEncoding.EncodeToString(b))
	} else {
		req, err = http.NewRequest("POST", c.buildHttpUrl(), bytes.NewReader(b))
		if err != nil {
			n = 0
			return
		}
	}

	req.Header.Set("User-Agent", defaultUserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		n = 0
		return
	}
	defer resp.Body.Close()

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
