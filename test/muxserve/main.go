package main

import (
	"fmt"
	"net/http"
	"strings"
)

func main() {
	m := http.NewServeMux()
	m.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("this is a hello")
	})
	m.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request.RequestURI)
		fmt.Println(len(strings.Split(request.RequestURI, "/")))
	})

	s := http.Server{}
	s.Handler = m
	s.Addr = ":8080"
	s.ListenAndServe()
}
