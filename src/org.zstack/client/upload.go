package client

import (
	"flag"
	"fmt"
	"org.zstack/lib"
	"path"
)

type (
	Upload struct {
	}
)

func (u *Upload) Name() string {
	return "upload"
}

func (u *Upload) Flags(f *flag.FlagSet) {
}

func (u *Upload) CheckFlags() error {
	args := flag.Args()[1:]
	if len(args) != 2 {
		return fmt.Errorf("Wrong parameters. Usage: zanker upload absolute_path_to_src_file absolute_path_to_dst_file")
	}

	src := args[0]
	if !lib.IsFile(src) {
		return fmt.Errorf("%s is not found or not a file", src)
	}

	dst := args[1]
	if !path.IsAbs(dst) {
		return fmt.Errorf("%s is not an absolute path")
	}

	return nil
}

func (u *Upload) Run() int {
	return 0
}
