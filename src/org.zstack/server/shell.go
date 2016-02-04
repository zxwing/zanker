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
	ShellRoute = &Route{
		Path:    ApiURL("/shell"),
		Handler: APIShell,
	}
)

type (
	_shellResult struct {
		Code   int    `json:"code"`
		Stdout string `json:"stdout"`
		Stderr string `json:"stderr"`
	}
)

func _shell(f string, v ...interface{}) (retCode int, stdout, stderr string, err error) {
	var so, se bytes.Buffer
	command := exec.Command("bash", "-c", fmt.Sprintf(f, v...))
	command.Stdout = &so
	command.Stderr = &se

	var waitStatus syscall.WaitStatus
	if err = command.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			retCode = waitStatus.ExitStatus()
		}
	} else {
		waitStatus = command.ProcessState.Sys().(syscall.WaitStatus)
		retCode = waitStatus.ExitStatus()
	}

	stdout = string(so.Bytes())
	stderr = string(se.Bytes())

	return
}

func APIShell(w http.ResponseWriter, req *http.Request) {
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

	code, stdout, stderr, err := _shell(command)
	if err != nil {
		panic(err)
	}

	ret := &_shellResult{
		Code:   code,
		Stdout: stdout,
		Stderr: stderr,
	}

	if code == 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
	}

	json.NewEncoder(w).Encode(ret)
}

func init() {
	ShellRoute.Register()
}
