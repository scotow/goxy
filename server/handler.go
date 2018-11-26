package server

import (
	"github.com/scotow/goxy/common"
	"io"
	"net/http"
)

type clientCreationHandler 	func(http.ResponseWriter)
type clientReadHandler 		func(*Conn, io.Writer)
type clientWriteHandler 	func(*Conn, io.Reader)

type handler struct {
	connM		*connMapping

	creationH 	clientCreationHandler
	readH 		clientReadHandler
	writeH		clientWriteHandler
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if

	if h.readH == nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}


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
