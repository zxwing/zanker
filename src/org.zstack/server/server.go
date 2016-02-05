package server

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	//"io/ioutil"
	"net"
	"net/http"
	"os"

	LOG "github.com/Sirupsen/logrus"
)

type (
	Server struct {
		listener net.Listener
		handler  http.Handler
		router   *mux.Router
	}
)

var (
	optionSocketPath string

	server = &Server{}
)

func (serv *Server) Run() {
	go func() {
		if err := http.Serve(serv.listener, serv.router); err != nil {
			panic(fmt.Errorf("cannot start the HTTP server, %v", err))
		}
	}()
}

func prepareRoutes() {
	server.router = NewRouter()
}

func prepareSocket() {
	if _, err := os.Stat(optionSocketPath); os.IsNotExist(err) {
		LOG.Debugf("the socket[%s] is not found, create a new one\n", optionSocketPath)
	}

	l, err := net.Listen("unix", optionSocketPath)
	if err != nil {
		panic(fmt.Errorf("cannot listen on the socket[%s], %v", optionSocketPath, err))
	}

	server.listener = l
}

func New() *Server {
	parseOptions()
	prepareRoutes()
	prepareSocket()

	return server
}

func registerServerOptions() {
	o := &Option{}

	o.BeforeParse = func() {
		flag.StringVar(&optionSocketPath, "socket", "", "path to the unix socket")
	}

	o.AfterParse = func() {
		if optionSocketPath == "" {
			flag.Usage()
			fmt.Printf("option [-socket] is required and cannot be an empty string\n")
			os.Exit(1)
		}
	}

	registerOption(o)
}

func initLog() {
	LOG.SetOutput(os.Stderr)
	LOG.SetLevel(LOG.DebugLevel)
}

func init() {
	initLog()
	registerServerOptions()
}
