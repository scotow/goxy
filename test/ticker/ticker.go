package main

import (
	"fmt"
	"os"
	"time"
)

func waitForStop(done chan bool) {
	b := make([]byte, 1)
	os.Stdin.Read(b)
	done <- true
}

func main() {
	doneC := make(chan bool)
	go waitForStop(doneC)

	t := time.NewTicker(time.Second)
	for {
		select {
		case <-t.C:
			fmt.Println("Waiting...")
		case <-doneC:
			return
		}
	}
}
