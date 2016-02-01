package main

import (
	"org.zstack/client"
	"os"
)

func main() {
	c := client.NewClient()
	c.Run()
	os.Exit(0)
}
