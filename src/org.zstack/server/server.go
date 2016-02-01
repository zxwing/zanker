package server

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

type (
	Server struct {
		listener net.Listener
		handler  http.Handler

		SocketPath string
	}
)

const (
	ECHO_PATH = "/echo"
)

func echoHandler(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Printf("body:%s\n", body)
	fmt.Fprint(w, "hello world!")
}

func (serv *Server) Run() {
	go func() {
		if err := http.Serve(serv.listener, nil); err != nil {
			panic(fmt.Errorf("cannot start the HTTP server, %v", err))
		}
	}()
}

func New() *Server {
	var socketPath string
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
		fmt.Printf("the socket[%s] is not found, create a new one\n", socketPath)
	}

	l, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(fmt.Errorf("cannot listen on the socket[%s], %v", socketPath, err))
	}

	serv := &Server{
		SocketPath: socketPath,
		listener:   l,
	}

	http.HandleFunc(ECHO_PATH, echoHandler)

	return serv
}
