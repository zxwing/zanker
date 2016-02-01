package main

import (
	"flag"
	"fmt"
	"org.zstack/server"
	"os"
)

var (
	socketPath string
)

func main() {
	flag.StringVar(&socketPath, "socket", "", "path to the unix socket")
	flag.Parse()

	if flag.NArg() > 0 {
		flag.Usage()
		fmt.Printf("unknown options %v\n", flag.Args())
		os.Exit(1)
	}

	if socketPath == "" {
		flag.Usage()
		fmt.Printf("option [-socket] is required and cannot be an empty string\n")
		os.Exit(1)
	}

	if _, err := os.Stat(socketPath); os.IsNotExist(err) {
		fmt.Printf("the socket[%s] is not found, create a new one", socketPath)
	}

	serv := server.New(socketPath)
	serv.Start()
	select {}
}
