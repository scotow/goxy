package main

import (
	"bytes"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net"
	"net/http"
)

var (
	buffer bytes.Buffer
	conn   *net.TCPConn
)

func initSsh() {
	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:22")
	if err != nil {
		log.Panic(err)
	}

	c, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		log.Panic(err)
	}
	conn = c

	b := make([]byte, 1024)
	for {
		n, _ := conn.Read(b)
		buffer.Write(b[:n])
	}
}

func handleOutput(w http.ResponseWriter, _ *http.Request) {
	log.Println("Handling output...")

	io.Copy(w, &buffer)
}

func handleInput(_ http.ResponseWriter, r *http.Request) {
	log.Println("Handling input...")

	defer r.Body.Close()
	io.Copy(conn, r.Body)
}

func main() {
	go initSsh()

	r := mux.NewRouter()
	r.HandleFunc("/", handleOutput).Methods("GET")
	r.HandleFunc("/", handleInput).Methods("POST")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", r))
}
