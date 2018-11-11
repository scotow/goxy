package server2

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"sync"
)

func NewListener(localAddr *net.TCPAddr) (*Listener, error) {
	l := Listener{}
	l.localAddr = localAddr

	r := mux.NewRouter()

	//r.HandleFunc("/status", l.handleStatus).Methods("GET")
	r.HandleFunc("/create", l.handleAccept).Methods("GET", "POST")
	//r.HandleFunc("/{id}/close", l.handleClose).Methods("GET", "POST")
	r.HandleFunc("/write/{id}", l.handleClientOutput).Methods("POST")
	r.HandleFunc("/read/{id}", l.handleClientFetch).Methods("POST")

	l.server = &http.Server{}
	l.server.Addr = localAddr.String()
	l.server.Handler = r

	l.acceptC = make(chan *Conn)
	l.connections = make(map[string]*Conn)

	return &l, nil
}

type Listener struct {
	localAddr *net.TCPAddr

	server  *http.Server
	acceptC chan *Conn

	cLock       sync.RWMutex
	connections map[string]*Conn
}

func (l *Listener) Start() error {
	return l.server.ListenAndServe()
}

func (l *Listener) getConnection(r *http.Request) (*Conn, string) {
	id := mux.Vars(r)["id"]

	l.cLock.RLock()
	defer l.cLock.RUnlock()

	return l.connections[id], id
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

func (l *Listener) handleAccept(w http.ResponseWriter, r *http.Request) {
	remoteAddr, err := net.ResolveTCPAddr("tcp", r.RemoteAddr)

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	conn := newConn(l.localAddr, remoteAddr)

	l.cLock.Lock()
	l.connections[conn.id] = conn
	l.cLock.Unlock()

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, conn.id)

	l.acceptC <- conn
}

func (l *Listener) handleClientOutput(w http.ResponseWriter, r *http.Request) {
	conn, _ := l.getConnection(r)

	if conn == nil {
		http.Error(w, "cannot find connection with id", http.StatusBadRequest)
		return
	}

	remaining := int(r.ContentLength)

	for {
		b := <-conn.readC
		n, err := r.Body.Read(b)
		conn.readNC <- n

		if err != nil {
			break
		}

		remaining -= n

		if remaining == 0 {
			break
		}
	}
}

func (l *Listener) handleClientFetch(w http.ResponseWriter, r *http.Request) {
	conn, _ := l.getConnection(r)

	if conn == nil {
		http.Error(w, "cannot find connection with id", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	max, err := strconv.Atoi(string(data))
	if err != nil {
		return
	}

	b := <-conn.writeC

	if len(b) < max {
		max = len(b)
	}

	n, _ := w.Write(b[:max])
	conn.writeNC <- n
}
