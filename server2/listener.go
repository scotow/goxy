package server2

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
)

const (
	payload = "PStvSlrIqM0b66tV6zblnEVLRXl0iniD14YNYvnkFw0A7N7MNWtSTcaZegogOULtJwPMAm3Bz7RSLA6oSBvS44YOWhfj3phLluc3mdsUJ0O25qzfldGeeTicZhPZqptVa4Nn0d/4j9h7vjiOI9W0MEDyYGWcSvCkXMwgpvrQ+uPBs9S3W2mtBaKckcNROtcWWzw57Yefi2x98cjEgEIij9SJugGx7xIE6sB+N0qg17qOSPDyVwwQvG+5bnwtD+krkVc63j1KzDvbkaDp+DNq+DQqj5IRqaEEbVoxReL+TXCykoA4l+vEl/fAd6wydZ5NUF5RNKoVSSASrfZehAowUxq+AqWhvldeCpH6F1DhJ3gl23nR80pkTs9QMY/lYOJyrB4zvlbps/QZbylbiJSNY0VaVPaF68KgD//GoItL07ka1e7p3peMr/fEKTyOTMNNWRdnVsdLPA0WJtgZNMFPMaiT3EgiDClWK6xpNRnf4l3JRB2335u7d0IAxter+Y1O+Buf+WvFamaQHSQkAHRhEjp/nqII37SCltplCTIH0UvJcND8rzoqcvOiLdojZ4y+DLFx2r6c/s6/T7SWTPso32hM6nQTlWtvnyoBgiGHZCaRPgwR3KTkXD71sC4qIax+XMJuI/60YQqHA9AJUxosDEpZjkT+AUca72U7Q5jckYi6otI22PJhFuo+FLtdPbA+CpAy/R/dp3WsQ3hQAtMBNQKrikwC1snTU8hmDQbwmEhDmIZzJIKiFdqWqBVvoZjC4wN5aQrjNdOtZHDWd554Vw4MlNIHMj4m+cMQhaNy3DseAWztCWq68A6u83SzIwghQiI2kuIhyRUB0JINo+fVQf5ASMkSXOlWEb85wU5DkyHN1ii9Of65rflWVM9pFBWK1vD1R9UcywiXPzVlvjE93Ged7BsTMSexEuFHVmx10XydqBMg58jJ0v05v6j56CtP/ftY03wLW2xhz8m1SVo3kIQt7FovtAlbd4M3p7Dr6cCOAPvWiUqlRuZGh/JcwEx7hYkmcV9RmnBNZlBeo+6i78reO9Koc6Lu3j4bD4Cp5CuKSod+6RygYCUH0adeI/UbDU3NXMLyKshkLn/VQ7CiFgPUkIskk9PwUm0HHtEZHNd3PzN606VlJmOjNki54nY9UIL1bYQul98C6UZTQI2oX2dKGoHdOW3/31CAT4nBhppG8aP9Lzy1C8U7MeUjeSFkmuoxuz6wf+rKZRT8Ml3b7V/SZeGoMr9JdoY18ftmk4Cir95nDIAxHqT91AvHVsqGTTmixRJ4E1iQeUTqCnH29Ay7rsmMQFuOOBFclICY6Mog8XeTMsACJzjv7Im/XZkL5jwUKxwbtcd7Ak0sXFQjq6oGCFtGzQHBnHvTcHDym11WjFJrFm9oPfB1zpLilQd0WZ7MqDP0RPicaohRoGbp9z5z2I5Lu3kQs6a2UFW37WDelkA8ac12JiNvhiSFGiSn6KkiFf/5dD9UwLc0zuytQYDkvrCraFLro5UfkjV8V/CuGbqBcIPgnoNH7qrIRvd+ZbICGZFJdd71m+1u8vowuseWvGBWA1JBW0oaXdQ60FcqwcxqBI6P0GB/b3TX6+giPAhDq0AP3x/FtY84CJnY9TAHsf+XYgjGIdLdkaSAReQCArvndB6QZYAPAXb9Y/vmQ4RD0whW60VdOOlg5c8y1qJtrZz8koavvIs9l7GcP8eC4zIleEugMU7s737ESEge8mIk2cSda/JZHB5IjNIQS7zub2qkOjU96EJ/OFm7wrjjNVFjp4ldW9d60Wh80lHuOFiuHpYWOcE6j0gFozKIF2KQ0FW1zED2b5R/zzELBzA/TfEtTb8iDa8kJZhrA2QQ7QkkN5AwLSPckox3MFEc8iUJOLUDr5GLkDH2kTWXQ/SsfMJ930DNqqkHijbQNdP9gEmNKOQgo2A="
	//payloadBackup = "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"
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

func (l *Listener) Start() {
	log.Panic(l.server.ListenAndServe())
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

	conn, err := newConn(l.localAddr, remoteAddr)
	if err != nil {
		errCode := http.StatusInternalServerError
		http.Error(w, http.StatusText(errCode), errCode)
		return
	}

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
		fmt.Println("cannot find connection with id")
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("cannot read fetch request body")
		return
	}

	clientCapacity, err := strconv.Atoi(string(data))
	if err != nil {
		fmt.Println("cannot parse  fetch request buffer size")
		return
	}

	b := <-conn.writeC

	var max int
	if clientCapacity < len(b) {
		max = clientCapacity
		fmt.Println("fetch: asked is too big, sending what we got")
	} else {
		max = len(b)
	}

	//payloadDecoded, err := base64.StdEncoding.DecodeString(payload)

	//w.Header().Set("X-Will-Write", fmt.Sprint(len(payloadDecoded)))
	w.Header().Set("X-Will-Write", fmt.Sprint(max))
	w.Header().Set("Content-Type", "application/octet-stream")
	//w.Header().Set("Content-Length", fmt.Sprint(max))

	//encoded := base64.StdEncoding.EncodeToString(b)
	//fmt.Fprintf(conn.logger, encoded)

	//n, err := w.Write(payloadDecoded)
	n, err := w.Write(b[:max])
	if err != nil || n != max {
		fmt.Println("error while writing content to client read request")
	}

	conn.writeNC <- n

	fmt.Fprintf(conn.logger, "http write: %d bytes wrote ouf of the %d asked ones\n", n, len(b))
}
