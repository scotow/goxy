package main

import (
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func newSshConn() (*net.TCPConn, error) {
	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:2222")
	if err != nil {
		return nil, err
	}

	sshConn, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		return nil, err
	}

	return sshConn, nil
}

func handleWsConn(w http.ResponseWriter, r *http.Request) {
	sshConn, err := newSshConn()
	if err != nil {
		log.Println("cannot open ssh connection", err)
		return
	}

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		for {
			messageType, reader, err := wsConn.NextReader()
			if err != err {
				log.Println(err)
				continue
			}

			log.Printf("Type: %d\n", messageType)
			io.Copy(sshConn, reader)
		}
	}()

	for {
		writer, err := wsConn.NextWriter(websocket.BinaryMessage)
		if err != err {
			log.Println(err)
			continue
		}

		log.Printf("New writter\n")
		io.Copy(writer, sshConn)
		writer.Close()
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "html/index.html")
}

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", handleWsConn)
	log.Fatal(http.ListenAndServe(":8080", nil))
}