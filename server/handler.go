package server

import (
	"encoding/base64"
	"github.com/pkg/errors"
	"github.com/scotow/goxy/common"
	"io"
	"net/http"
	"strings"
)

var (
	errNoHandlerFound = errors.New("cannot find suitable path for the request")
)

type clientCreationHandler func(http.ResponseWriter, string)
type clientReadHandler func(*Conn, http.ResponseWriter)
type clientWriteHandler func(*Conn, io.Reader)

func newHandler(mapping *connMapping) *handler {
	h := new(handler)
	h.connM = mapping
	return h
}

type handler struct {
	connM *connMapping

	creationH clientCreationHandler
	readH     clientReadHandler
	writeH    clientWriteHandler
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" {
		if h.creationH == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		h.creationH(w, r.RemoteAddr)
		return
	}

	if r.Method == "GET" && r.Header.Get("Authorization") != "" {
		if h.writeH == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		conn := h.findConnection(w, r)
		if conn == nil {
			return
		}

		h.writeH(conn, base64.NewDecoder(base64.StdEncoding, strings.NewReader(r.Header.Get("Authorization"))))
		return
	}

	if r.Method == "POST" {
		if h.writeH == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		conn := h.findConnection(w, r)
		if conn == nil {
			return
		}

		h.writeH(conn, r.Body)
		return
	}

	if r.Method == "GET" {
		if h.readH == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		conn := h.findConnection(w, r)
		if conn == nil {
			return
		}

		h.readH(conn, w)
		return
	}

	http.Error(w, errNoHandlerFound.Error(), http.StatusInternalServerError)
}

func (h *handler) findConnection(w http.ResponseWriter, r *http.Request) (conn *Conn) {
	if h.connM == nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	token := common.TokenFromPath(r.RequestURI[1:])
	conn, err := h.connM.getConn(token)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	return
}
