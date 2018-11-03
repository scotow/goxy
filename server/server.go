package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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
	s.router.HandleFunc("/create", s.create).Methods("POST")
	s.router.PathPrefix("/").HandlerFunc(s.input).Methods("POST")
	s.router.PathPrefix("/").HandlerFunc(s.output).Methods("GET", "POST")

	return &s, nil
}

func (s *Server) status(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
	log.WithField("origin", r.RemoteAddr).Info("Received status request.")
}

func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid creation request", http.StatusBadRequest)
		log.WithField("error", err).Warn("Invalid creation request body.")
		return
	}

	var version, remoteAddr string
	reqString := string(reqBytes)
	n, err := fmt.Sscanf(reqString, "GOXY %s %s", &version, &remoteAddr)

	if n != 2 || err != nil {
		http.Error(w, "invalid creation request", http.StatusBadRequest)
		log.WithField("header", reqString).Warn("Invalid creation request header.")
		return
	}

	id := strconv.FormatInt(time.Now().UnixNano(), 10)

	conn, err := newConnection(remoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.WithField("error", err.Error()).Warn("Cannot create connection.")
		return
	}

	s.connections[id] = conn
	go conn.waitForOutput()

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, id)
	log.WithFields(log.Fields{
		"id":         id,
		"origin":     r.RemoteAddr,
		"remoteAddr": remoteAddr,
	}).Info("Connection created.")
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
	log.WithField("address", s.httpAddr).Info("Starting HTTP server.")
	return http.ListenAndServe(s.httpAddr, s.router)
}
