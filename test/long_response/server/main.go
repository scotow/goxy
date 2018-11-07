package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func handleSlow(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "sending first line of data")
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	} else {
		log.Println("Damn, no flush")
	}
	time.Sleep(3 * time.Second)
	fmt.Fprintln(w, "sending second line of data")
}

func handleFast(w http.ResponseWriter, r *http.Request) {
	io.Copy(os.Stdout, r.Body)
	fmt.Fprintln(w, "sending first line of data")
	fmt.Fprintln(w, "sending second line of data")
}

func main() {
	http.HandleFunc("/", handleFast)
	http.HandleFunc("/slow", handleSlow)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
