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

	s.router.HandleFunc("/status", s.handleStatus).Methods("GET")
	s.router.HandleFunc("/create", s.handleCreate).Methods("POST")
	s.router.HandleFunc("/{id}/close", s.handleClose).Methods("GET", "POST")
	s.router.HandleFunc("/{id}", s.handleClientOutput).Methods("POST")
	s.router.HandleFunc("/{id}", s.handleClientFetch).Methods("GET", "POST")

	return &s, nil
}

func (s *Server) Start() error {
	log.WithField("address", s.httpAddr).Info("Starting HTTP server.")
	return http.ListenAndServe(s.httpAddr, s.router)
}

// HTTP Handlers

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
	log.WithField("origin", r.RemoteAddr).Info("Received status request.")
}

func (s *Server) handleCreate(w http.ResponseWriter, r *http.Request) {
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

	go conn.start()

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, id)
	log.WithFields(log.Fields{
		"id":         id,
		"origin":     r.RemoteAddr,
		"remoteAddr": remoteAddr,
	}).Info("Connection created.")
}

func (s *Server) handleClientOutput(w http.ResponseWriter, r *http.Request) {
	s.handleIfConnected(w, r, s.clientOutput)
}

func (s *Server) handleClientFetch(w http.ResponseWriter, r *http.Request) {
	s.handleIfConnected(w, r, s.clientFetch)
}

func (s *Server) handleClose(w http.ResponseWriter, r *http.Request) {
	s.handleIfConnected(w, r, s.clientClose)
}

func (s *Server) handleIfConnected(w http.ResponseWriter, r *http.Request, handler connectionHandler) {
	conn, id := s.getConnection(r)
	if conn == nil {
		http.Error(w, "invalid connection id", http.StatusBadRequest)
		log.WithFields(log.Fields{
			"id":      id,
			"address": r.RemoteAddr,
		}).Warn("Invalid connection id.")
		return
	}

	handler(conn, id, w, r)
}

// Connection methods.

type connectionHandler func(*connection, string, http.ResponseWriter, *http.Request)

func (s *Server) getConnection(r *http.Request) (*connection, string) {
	id := mux.Vars(r)["id"]

	s.connectionsLock.RLock()
	defer s.connectionsLock.RUnlock()

	return s.connections[id], id
}

func (s *Server) clientOutput(c *connection, id string, w http.ResponseWriter, r *http.Request) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.closing {
		w.WriteHeader(http.StatusGone)
		s.deleteConnection(id)
		return
	}

	io.Copy(c.tcpConn, r.Body)
	r.Body.Close()
}

func (s *Server) clientFetch(c *connection, id string, w http.ResponseWriter, _ *http.Request) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.closing {
		w.WriteHeader(http.StatusGone)
		s.deleteConnection(id)
	}

	io.Copy(w, &c.outputBuffer)
}

func (s *Server) clientClose(c *connection, id string, w http.ResponseWriter, r *http.Request) {
	s.deleteConnection(id)
	c.clientClosed <- true

	w.WriteHeader(http.StatusOK)
	log.WithFields(log.Fields{
		"id":      id,
		"address": r.RemoteAddr,
	}).Info("Connection closed.")
}

func (s *Server) deleteConnection(id string) {
	s.connectionsLock.Lock()
	defer s.connectionsLock.Unlock()

	delete(s.connections, id)
}
