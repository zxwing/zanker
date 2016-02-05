package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"syscall"
)

var (
	PluginShell = &Shell{}
)

type (
	_shellResult struct {
		Code   int    `json:"code"`
		Stdout string `json:"stdout"`
		Stderr string `json:"stderr"`
	}

	_shell struct{}

	Shell _shell
)

func (s *Shell) Methods() []string {
	return []string{"GET", "POST"}
}

func (s *Shell) Path() []string {
	return []string{"/shell"}
}

func (s *Shell) Handler(w http.ResponseWriter, req *http.Request) {
	var command string
	command = req.URL.Query().Get("command")
	if command == "" {
		var buf bytes.Buffer
		buf.ReadFrom(req.Body)
		command = buf.String()
	}

	if command == "" {
		panic(fmt.Errorf("command cannot be empty. Please either set the command by a query string, for example, ?command=ls or put the command in the body"))
	}

	code, stdout, stderr := s.shell(command)

	ret := &_shellResult{
		Code:   code,
		Stdout: stdout,
		Stderr: stderr,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)

}

func (s *Shell) shell(f string, v ...interface{}) (retCode int, stdout, stderr string) {
	var so, se bytes.Buffer
	command := exec.Command("bash", "-c", fmt.Sprintf(f, v...))
	command.Stdout = &so
	command.Stderr = &se

	var waitStatus syscall.WaitStatus
	if err := command.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			retCode = waitStatus.ExitStatus()
		} else {
			// looks like a system error, for example, IO error
			panic(err)
		}
	} else {
		waitStatus = command.ProcessState.Sys().(syscall.WaitStatus)
		retCode = waitStatus.ExitStatus()
	}

	stdout = string(so.Bytes())
	stderr = string(se.Bytes())

	return
}

func init() {
	RegisterApiRoute(PluginShell)
}
