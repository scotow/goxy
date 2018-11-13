package client2

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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

	logFile, err := os.Create(fmt.Sprintf("goxy-client%s.log", id))
	if err != nil {
		return nil, err
	}

	logger := io.MultiWriter(os.Stdout, logFile)

	conn := Conn{string(id), remoteAddr, logFile, logger}
	return &conn, nil
}

type Conn struct {
	id         string
	remoteAddr *net.TCPAddr

	logFile *os.File
	logger  io.Writer
}

func (c *Conn) Read(b []byte) (n int, err error) {
	httpAddr := fmt.Sprintf("http://%s/read/%s", c.remoteAddr.String(), c.id)

	resp, err := http.Post(httpAddr, "*/*", strings.NewReader(strconv.Itoa(len(b))))
	if err != nil {
		n = 0
		fmt.Println("HTTP read POST request", err.Error())
		return
	}
	defer resp.Body.Close()

	n, err = resp.Body.Read(b)
	if err != nil && err != io.EOF {
		fmt.Printf("Read: reading from body response error: %s\n", err.Error())
		return
	}

	writtenServerS := resp.Header.Get("X-Will-Write")
	if writtenServerS == "" {
		fmt.Println()
		fmt.Println("missing X-Will-Write header from response")
		return
	} else {
		fmt.Printf("X-Will-Write: %s\n", writtenServerS)
	}

	writtenServer, err := strconv.Atoi(writtenServerS)
	if err != nil {
		fmt.Println("cannot parse X-Will-Write header")
	} else {
		if writtenServer != n {
			fmt.Println("miss matching written/read length")
			fmt.Println("http content length: ", resp.ContentLength)
			fmt.Fprintf(c.logger, "Read: buffer size: %d. Read: %d.\n", len(b), n)

			fmt.Println(base64.StdEncoding.EncodeToString(b[:n]))
			//fmt.Println(string(b[:n]))
			return
		}
	}

	// TODO: Check for error on read (should be EOF).
	fmt.Fprintf(c.logger, "Read: buffer size: %d. Read: %d.\n", len(b), n)
	err = nil
	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	httpAddr := fmt.Sprintf("http://%s/write/%s", c.remoteAddr.String(), c.id)

	resp, err := http.Post(httpAddr, "*/*", bytes.NewReader(b))
	if err != nil {
		n = 0
		return
	}
	defer resp.Body.Close()

	n = len(b)
	// TODO: Check for end of file with custom HTTP status code.

	fmt.Fprintf(c.logger, "Write: buffer size: %d. Written: %d.\n", len(b), n)
	return
}

func (c *Conn) Close() error {
	c.logFile.Close()
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
