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
type clientReadHandler func(*Conn, http.ResponseWriter, *common.Hider)
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

		if conn.writeToken != r.Header.Get("Referer") {
			w.Write([]byte("The page you asked was deleted."))
			return
		}

		conn.newWriteToken()
		w.Header().Set("X-Referer", conn.writeToken)

		h.writeH(conn, base64.NewDecoder(base64.StdEncoding, strings.NewReader(r.Header.Get("Authorization"))))

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("The page you asked was deleted."))
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

		if conn.writeToken != r.Header.Get("Referer") {
			w.Write([]byte("The page you asked was deleted."))
			return
		}

		conn.newWriteToken()
		w.Header().Set("X-Referer", conn.writeToken)

		h.writeH(conn, r.Body)

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("The page you asked was renamed."))
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

		hider, err := common.HiderFromPath(r.RequestURI)
		if err != nil {
			w.Write([]byte("The page you asked was deleted."))
			//http.Error(w, http.StatusText(http.StatusForbidden), http.StatusInternalServerError)
			return
		}

		if conn.readToken != r.Header.Get("Referer") {
			w.Write([]byte("The page you asked was deleted."))
			return
		}

		conn.newReadToken()
		w.Header().Set("X-Referer", conn.readToken)

		w.Header().Set("Content-Type", hider.Mime)
		h.readH(conn, w, hider)
		return
	}

	http.Error(w, errNoHandlerFound.Error(), http.StatusInternalServerError)
}

func (h *handler) findConnection(w http.ResponseWriter, r *http.Request) (conn *Conn) {
	if h.connM == nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	parts := strings.SplitN(r.RequestURI[1:], ".", 2)
	token := common.TokenFromPath(parts[0])
	conn, err := h.connM.getConn(token)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	return
}
