package server

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

type (
	Server struct {
		listener net.Listener
		handler  http.Handler
		//mux      *http.ServeMux

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

func (serv *Server) Start() {
	go func() {
		if err := http.Serve(serv.listener, nil); err != nil {
			panic(fmt.Errorf("cannot start the HTTP server, %v", err))
		}
	}()
}

func New(socketPath string) *Server {
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
