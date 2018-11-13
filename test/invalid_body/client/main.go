package main

import (
	"log"
	"net/http"
)

func main() {
	resp, err := http.Post("http://public.tiir.local:8081", "*/*", nil)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	b := make([]byte, 3000)
	n, err := resp.Body.Read(b)

	log.Println(n)
}
