package main

import (
	"time"

	goxy "github.com/scotow/goxy/client"
)

func main() {
	client, _ := goxy.NewClient(2222, "localhost:80", "localhost:22")
	client.WaitUntilServerUp(time.Second)
	client.Start()
}
