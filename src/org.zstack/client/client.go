package client

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	//"io/ioutil"
	"net"
	"net/http"
	"os"
)

type (
	Client struct {
		conn   net.Conn
		client *http.Client

		SocketPath string
	}
)

var (
	//test       string
	socketPath string
)

const (
	SCHEME = "unix+http"
)

func (client *Client) RoundTrip(req *http.Request) (*http.Response, error) {
	var buf bytes.Buffer
	req.Write(&buf)
	if _, err := client.conn.Write(buf.Bytes()); err != nil {
		return nil, err
	}

	rsp, err := http.ReadResponse(bufio.NewReader(client.conn), req)
	return rsp, err
}

func NewClient() *Client {
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
		fmt.Printf("the socket[%s] not found\n", socketPath)
		os.Exit(1)
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		panic(fmt.Errorf("cannot connect to the socket[%s]", socketPath))
	}

	c := &Client{
		SocketPath: socketPath,
		conn:       conn,
	}

	t := &http.Transport{}
	t.RegisterProtocol(SCHEME, c)
	c.client = &http.Client{Transport: t}

	return c
}

func (client *Client) Run() {
	rsp, err := client.client.Get(fmt.Sprintf("%s:///v1/shell?command=ls", SCHEME))
	if err != nil {
		panic(err)
	}

	defer rsp.Body.Close()
	var res bytes.Buffer
	res.ReadFrom(rsp.Body)
	fmt.Println(res.String())
}
