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
	buffer  bytes.Buffer
	sshConn *net.TCPConn
)

func initSsh() {
	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:5555")
	if err != nil {
		log.Panic(err)
	}

	conn, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		log.Panic(err)
	}

	sshConn = conn
	io.Copy(&buffer, sshConn)
}

func handleOutput(w http.ResponseWriter, r *http.Request) {
	io.Copy(w, &buffer)
}

func handleInput(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling input....")
	defer r.Body.Close()
	io.Copy(sshConn, r.Body)
}

func main() {
	// Listen on /comment for HTTP post
	// Send data received on /comment to stdout
	// On each HTTP get request on /video sends buffer

	go initSsh()

	r := mux.NewRouter()
	r.HandleFunc("/", handleOutput).Methods("GET")
	r.HandleFunc("/", handleInput).Methods("POST")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", r))
}
