package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var (
	count int
)

func main() {
	for i := 0; i < 100; i++ {
		resp, err := http.Post("http://public.tiir.local:8081", "*/*", nil)
		if err != nil {
			log.Println(err)
		}

		//b := make([]byte, 5000)
		//n, err := resp.Body.Read(b)

		b, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		n := len(b)

		readHeader := resp.Header.Get("X-GOXY")
		read, _ := strconv.Atoi(readHeader)

		if n != read {
			log.Printf("%d: Miss matching size: n = %d, w = %d\n", count, n, read)
		}

		count++
	}
}
