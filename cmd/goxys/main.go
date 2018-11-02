package main

import (
	"flag"
	goxy "github.com/scotow/goxy/server"
)

var (
	localHTTP = flag.String("p", ":8080", "local HTTP address used by goxy clients (address:port)")
)

func main() {
	flag.Parse()

	server, _ := goxy.NewServer(*localHTTP)
	server.Start()
}
