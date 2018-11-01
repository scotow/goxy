package main

import (
	goxy "github.com/scotow/goxy/server"
)

func main() {
	server, _ := goxy.NewServer(":80")
	server.Start()
}
