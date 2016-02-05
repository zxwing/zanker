package client

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	LIB "org.zstack/lib"
	"os"
)

type (
	_shell struct {
		command string
		file    string
		json    bool
		stdin   bool
	}

	_shellResult struct {
		Code   int    `json:"code"`
		Stdout string `json:"stdout"`
		Stderr string `json:"stderr"`
	}
)

const (
	SHELL_API_PATH = "/shell"
)

func (s *_shell) Name() string {
	return "shell"
}

func (s *_shell) Flags(f *flag.FlagSet) {
	f.StringVar(&s.command, "c", "", "one line shell command")
	f.StringVar(&s.file, "f", "", "path to the shell file")
	f.BoolVar(&s.json, "json", false, "encode output in JSON format")
	f.BoolVar(&s.stdin, "stdin", false, "read the script from the stdin")
}

func (s *_shell) CheckFlags() error {
	if s.command == "" && s.file == "" && !s.stdin {
		return fmt.Errorf("please specify a shell command by '-c', a shell script file by '-f', or import a script by the '-stdin'")
	} else if s.command == "" && s.file != "" {
		if !LIB.IsFile(s.file) {
			return fmt.Errorf("the shell file[%s] not found or not a file\n", s.file)
		}
	}

	return nil
}

func (s *_shell) Run() int {
	var status int
	var body string
	var rsp *http.Response

	if s.command != "" {
		url := Url(SHELL_API_PATH).Query("command", s.command).String()
		status, body, rsp = Http.Get(url)
	} else if s.file != "" {
		content, err := ioutil.ReadFile(s.file)
		if err != nil {
			panic(err)
		}

		url := Url(SHELL_API_PATH).String()
		status, body, rsp = Http.Post(url, string(content))
	} else if s.stdin {
		content, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}

		url := Url(SHELL_API_PATH).String()
		status, body, rsp = Http.Post(url, string(content))
	}

	if status != 200 {
		panic(ServerError(fmt.Sprintf("%s, %s", rsp.Status, body)))
	}

	if s.json {
		fmt.Fprint(os.Stdout, body)
		return 0
	}

	ret := &_shellResult{}
	if err := json.NewDecoder(bytes.NewBufferString(body)).Decode(ret); err != nil {
		panic(err)
	}

	if ret.Stdout != "" {
		fmt.Fprint(os.Stdout, ret.Stdout)
	}
	if ret.Stderr != "" {
		fmt.Fprint(os.Stderr, ret.Stderr)
	}

	return ret.Code
}

func init() {
	RegisterSubcommand(&_shell{})
}
