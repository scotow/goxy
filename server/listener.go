package server

import (
	"encoding/base64"
	"fmt"
	"github.com/scotow/goxy/common"
	"io"
	"log"
	"net"
	"net/http"
)

var (
//ErrWriteLengthMismatch = errors.New("write and expected right mismatched")
)

func NewListener(localAddr *net.TCPAddr) (*Listener, error) {
	l := Listener{}
	l.localAddr = localAddr

	l.mapping = newConnMapping()

	handler := newHandler(l.mapping)
	handler.creationH = l.handleAccept
	handler.readH = l.handleClientFetch
	handler.writeH = l.handleClientOutput

	l.server = &http.Server{}
	l.server.Addr = localAddr.String()
	l.server.Handler = handler

	l.acceptC = make(chan *Conn)

	return &l, nil
}

type Listener struct {
	localAddr *net.TCPAddr

	mapping *connMapping
	server  *http.Server
	acceptC chan *Conn
}

func (l *Listener) Start() {
	log.Panic(l.server.ListenAndServe())
}

// Listener interface

func (l *Listener) Accept() (net.Conn, error) {
	return <-l.acceptC, nil
}

func (l *Listener) Close() error {
	panic("implement me")
}

func (l *Listener) Addr() net.Addr {
	return l.localAddr
}

// HTTP handlers

func (l *Listener) handleAccept(w http.ResponseWriter, rAddr string) {
	remoteAddr, err := net.ResolveTCPAddr("tcp", rAddr)

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	conn, err := newConn(l.localAddr, remoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = l.mapping.addConn(conn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	l.acceptC <- conn

	w.Header().Set("X-Referer", fmt.Sprintf("%s %s", conn.readToken, conn.writeToken))
	fmt.Fprintf(w, conn.id.Token())
}

func (l *Listener) handleClientOutput(conn *Conn, r io.Reader) {
	for {
		b := <-conn.readC
		n, err := r.Read(b)

		conn.readNC <- n

		if err == io.EOF {
			conn.readEC <- nil
			break
		}

		if err != nil {
			conn.readEC <- err
			break
		}

		conn.readEC <- nil
	}
}

func (l *Listener) handleClientFetch(conn *Conn, w http.ResponseWriter, hider *common.Hider) {
	b := <-conn.writeC

	n, err := w.Write(hider.HideData([]byte(base64.StdEncoding.EncodeToString(b))))
	if err != nil {
		fmt.Println("error while writing content to client read request")
	}

	if n == len(b) {
		fmt.Println("invalid hidden content length")
	}

	conn.writeNC <- len(b)
	conn.writeEC <- err
}
