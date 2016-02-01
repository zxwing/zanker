package main

import (
	"org.zstack/server"
)

func main() {
	serv := server.New()
	serv.Run()
	select {}
}
