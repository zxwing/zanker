package client

import (
	"flag"
	"fmt"
	"os"
)

type (
	_shell struct {
		command string
		file    string
	}
)

const (
	SHELL_API_PATH         = "/shell"
	SHELL_API_COMMAND_PATH = "/shell?command=%s"
)

func (s *_shell) Name() string {
	return "shell"
}

func (s *_shell) Flags(f *flag.FlagSet) {
	f.StringVar(&s.command, "c", "", "one line shell command")
	f.StringVar(&s.command, "f", "", "path to the shell file")
}

func (s *_shell) CheckFlags() error {
	if s.command == "" && s.file == "" {
		return fmt.Errorf("please specify either a shell command by '-c' or a shell script file by '-f'")
	} else if s.command == "" && s.file != "" {
		info, err := os.Stat(s.file)
		if os.IsNotExist(err) {
			return fmt.Errorf("the shell file[%s] not found\n", s.file)
		}

		if info.IsDir() {
			return fmt.Errorf("[%s] is a directory, not a file\n", s.file)
		}
	}

	return nil
}

func (s *_shell) Run() int {
	var status int
	var body string
	if s.command != "" {
		fmt.Printf("xxxxxxxxxxxxxxxxxxx %s\n", s.command)
		url := ApiURL(SHELL_API_COMMAND_PATH, s.command)
		fmt.Printf("yyyyyyyyyyyyyyyyy %s\n", url)
		status, body, _ = Http.Get(url)
	} else {
	}

	fmt.Println(body)
	if status == 200 {
		return 0
	} else {
		return 1
	}
}

func init() {
	RegisterSubcommand(&_shell{})
}
