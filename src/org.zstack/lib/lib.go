package lib

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func IsFile(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}

func Shell(f string, v ...interface{}) (retCode int, stdout, stderr string) {
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

func Md5(path string) string {
	code, stdout, stderr := Shell("/usr/bin/md5sum %s", path)
	if code != 0 {
		panic(fmt.Errorf("cannot calculate MD5 of the file[%s], %v, %s", path, code, stderr))
	}

	p := strings.Split(stdout, " +")
	return p[0]
}

func RandomString() string {
	f, err := os.Open("/dev/urandom")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	b := make([]byte, 32)
	f.Read(b)
	return string(b)
}
