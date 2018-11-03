package main

import (
	"flag"
	"time"

	goxy "github.com/scotow/goxy/client"
)

var (
	localTCP   = flag.String("l", ":2222", "local listening TCP address (host:port)")
	remoteHTTP = flag.String("h", "localhost:8080", "remote Goxy server HTTP address (host:port)")
	remoteTCP  = flag.String("r", "localhost:22", "remote TCP address (host:port)")
)

func main() {
	flag.Parse()

	client, _ := goxy.NewClient(*localTCP, *remoteHTTP, *remoteTCP)
	client.WaitUntilServerUp(time.Second)
	client.Start()
}
