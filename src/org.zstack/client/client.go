package client

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	//"io/ioutil"
	LOG "github.com/Sirupsen/logrus"
	"net"
	"net/http"
	"os"
)

type HttpClient struct {
	client *http.Client
}

var (
	_socketPath string

	Http *HttpClient
)

const (
	SCHEME = "unix+http"
)

type RoundTripFunc func(*http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func _prepareSocket() {
	if _socketPath == "" {
		flag.Usage()
		fmt.Printf("option [-socket] is required and cannot be an empty string\n")
		os.Exit(1)
	}

	if _, err := os.Stat(_socketPath); os.IsNotExist(err) {
		fmt.Printf("the socket[%s] not found\n", _socketPath)
		os.Exit(1)
	}

	conn, err := net.Dial("unix", _socketPath)
	if err != nil {
		panic(fmt.Errorf("cannot connect to the socket[%s]", _socketPath))
	}

	var rt RoundTripFunc = func(req *http.Request) (*http.Response, error) {
		var buf bytes.Buffer
		req.Write(&buf)
		if _, err := conn.Write(buf.Bytes()); err != nil {
			return nil, err
		}

		rsp, err := http.ReadResponse(bufio.NewReader(conn), req)
		return rsp, err
	}

	t := &http.Transport{}
	t.RegisterProtocol(SCHEME, rt)

	Http = &HttpClient{}
	Http.client = &http.Client{Transport: t}
}

func (h *HttpClient) Get(url string) (int, string, *http.Response) {
	LOG.Debugf("GET %s", url)
	rsp, err := h.client.Get(url)
	if err != nil {
		panic(err)
	}

	defer rsp.Body.Close()
	var res bytes.Buffer
	res.ReadFrom(rsp.Body)

	return rsp.StatusCode, res.String(), rsp
}

func (h *HttpClient) Post(url, body string) (int, string, *http.Response) {
	LOG.Debugf("POST %s", url)
	rsp, err := h.client.Post(url, "application/json", bytes.NewBufferString(body))

	if err != nil {
		panic(err)
	}

	defer rsp.Body.Close()
	var res bytes.Buffer
	res.ReadFrom(rsp.Body)

	return rsp.StatusCode, res.String(), rsp
}

func Run() {
	ParseSubCommands()
	_prepareSocket()
	RunSubCommand()
}

func init() {
	LOG.SetOutput(os.Stderr)
	LOG.SetLevel(LOG.DebugLevel)

	flag.StringVar(&_socketPath, "socket", "", "path to the unix socket")
}
