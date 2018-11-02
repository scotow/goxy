package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	httpAddr    string
	router      *mux.Router
	connections map[string]*connection
}

func NewServer(httpAddr string) (*Server, error) {
	s := Server{}
	s.httpAddr = httpAddr
	s.router = mux.NewRouter()
	s.connections = make(map[string]*connection)

	s.router.HandleFunc("/status", s.status).Methods("GET")
	s.router.HandleFunc("/create", s.create).Methods("GET", "POST")
	s.router.PathPrefix("/").HandlerFunc(s.input).Methods("POST")
	s.router.PathPrefix("/").HandlerFunc(s.output).Methods("GET")

	return &s, nil
}

func (s *Server) status(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	id := strconv.FormatInt(time.Now().UnixNano(), 10)

	conn, err := newConnection()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.connections[id] = conn
	go conn.waitForOutput()

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, id)
}

func (s *Server) input(w http.ResponseWriter, r *http.Request) {
	//log.Println("Input id:", r.Header.Get("X-Id"))
	id := r.RequestURI[1:]
	//log.Println("Input id:", id)

	conn := s.connections[id]

	defer r.Body.Close()
	io.Copy(conn.tcpConn, r.Body)
}

func (s *Server) output(w http.ResponseWriter, r *http.Request) {
	//log.Println("Output id:", r.Header.Get("X-Id"))
	id := r.RequestURI[1:]
	//log.Println("Output id:", id)

	conn := s.connections[id]

	conn.bufferLock.Lock()
	io.Copy(w, &conn.outputBuffer)
	conn.bufferLock.Unlock()
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.httpAddr, s.router)
}
