package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	r, _ := http.Get("http://localhost:8080/slow")
	io.Copy(os.Stdout, r.Body)
}
