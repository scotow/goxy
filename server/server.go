package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/scotow/goxy/common"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	httpAddr string
	router   *mux.Router

	connectionsLock sync.RWMutex
	connections     map[string]*connection
}

func NewServer(httpAddr string) (*Server, error) {
	s := Server{}
	s.httpAddr = httpAddr
	s.router = mux.NewRouter()
	s.connections = make(map[string]*connection)

	s.router.HandleFunc("/status", s.status).Methods("GET")
	s.router.HandleFunc("/create", s.create).Methods("POST")
	s.router.HandleFunc("/{id}/close", s.close).Methods("GET", "POST")
	s.router.HandleFunc("/{id}", s.outputClient).Methods("POST")
	s.router.HandleFunc("/{id}", s.outputServer).Methods("GET", "POST")

	return &s, nil
}

func (s *Server) status(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
	log.WithField("origin", r.RemoteAddr).Info("Received status request.")
}

func (s *Server) getConnection(r *http.Request) (*connection, string) {
	id := mux.Vars(r)["id"]

	s.connectionsLock.RLock()
	defer s.connectionsLock.RUnlock()

	return s.connections[id], id
}

func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid creation request", http.StatusBadRequest)
		log.WithField("error", err).Warn("Invalid creation request body.")
		return
	}

	var appName, version, remoteAddr string
	reqString := string(reqBytes)
	n, err := fmt.Sscanf(reqString, "%s %s %s", &appName, &version, &remoteAddr)

	if n != 3 || err != nil {
		http.Error(w, "invalid creation request", http.StatusBadRequest)
		log.WithField("header", reqString).Warn("Invalid creation request header.")
		return
	}

	if appName != common.AppName || version != common.Version {
		http.Error(w, "invalid app name or version", http.StatusBadRequest)
		log.WithField("header", reqString).Warn("Invalid app name or version.")
		return
	}

	id := strconv.FormatInt(time.Now().UnixNano(), 10)

	conn, err := newConnection(remoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.WithField("error", err.Error()).Warn("Cannot create connection.")
		return
	}

	// Add connections to map.
	s.connectionsLock.Lock()
	s.connections[id] = conn
	s.connectionsLock.Unlock()

	go s.startConnection(conn)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, id)
	log.WithFields(log.Fields{
		"id":         id,
		"origin":     r.RemoteAddr,
		"remoteAddr": remoteAddr,
	}).Info("Connection created.")
}

func (s *Server) outputClient(w http.ResponseWriter, r *http.Request) {
	conn, id := s.getConnection(r)
	if conn == nil {
		http.Error(w, "invalid connection id", http.StatusBadRequest)
		log.WithFields(log.Fields{
			"id":      id,
			"address": r.RemoteAddr,
		}).Warn("Invalid connection id.")
		return
	}

	if conn.closed {
		http.Error(w, "connection closed", http.StatusGone)
		s.deleteConnection(id)
		return
	}

	defer r.Body.Close()
	io.Copy(conn.tcpConn, r.Body)
}

func (s *Server) outputServer(w http.ResponseWriter, r *http.Request) {
	conn, id := s.getConnection(r)
	if conn == nil {
		http.Error(w, "invalid connection id", http.StatusBadRequest)
		log.WithFields(log.Fields{
			"id":      id,
			"address": r.RemoteAddr,
		}).Warn("Invalid connection id.")
		return
	}

	conn.bufferLock.Lock()
	io.Copy(w, &conn.outputBuffer)
	conn.bufferLock.Unlock()

	if conn.closed {
		http.Error(w, "connection closed", http.StatusGone)
		s.deleteConnection(id)
		return
	}
}

func (s *Server) close(w http.ResponseWriter, r *http.Request) {
	conn, id := s.getConnection(r)
	log.Println("found close request")
	if conn == nil {
		http.Error(w, "invalid connection id", http.StatusBadRequest)
		log.WithFields(log.Fields{
			"id":      id,
			"address": r.RemoteAddr,
		}).Warn("Invalid connection id.")
		return
	}

	s.deleteConnection(id)
	conn.stopBuffering <- true
}

func (s *Server) startConnection(conn *connection) {
	go conn.waitForOutput()

	<-conn.socketClosed
	conn.closed = true
}

func (s *Server) deleteConnection(id string) {
	s.connectionsLock.Lock()
	defer s.connectionsLock.Unlock()

	delete(s.connections, id)
}

func (s *Server) Start() error {
	log.WithField("address", s.httpAddr).Info("Starting HTTP server.")
	return http.ListenAndServe(s.httpAddr, s.router)
}
