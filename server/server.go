package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
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
	s.router.HandleFunc("/", s.input).Methods("POST")
	s.router.HandleFunc("/", s.output).Methods("GET")

	return &s, nil
}

func (s *Server) status(w http.ResponseWriter, r *http.Request) {
	log.Println(s.httpAddr)
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
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, id)
}

func (s *Server) input(w http.ResponseWriter, r *http.Request) {
	log.Println("Input id:", r.Header.Get("X-Id"))
}

func (s *Server) output(w http.ResponseWriter, r *http.Request) {
	log.Println("Output id:", r.Header.Get("X-Id"))
	//defer r.Body.Close()
	//io.Copy(conn, r.Body)
}

func (s *Server) Start() error {
	//http.Handle("/", s.router)
	return http.ListenAndServe(s.httpAddr, s.router)
}
