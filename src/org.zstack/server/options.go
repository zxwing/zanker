package server

import (
	"flag"
	"fmt"
	"os"
)

type (
	Option struct {
		BeforeParse func()
		AfterParse  func()
	}
)

var (
	options []*Option = make([]*Option, 0)
)

func registerOption(o *Option) {
	options = append(options, o)
}

func parseOptions() {
	for _, o := range options {
		o.BeforeParse()
	}

	flag.Parse()

	if flag.NArg() > 0 {
		flag.Usage()
		fmt.Printf("unknown option: %v\n", flag.Args())
		os.Exit(1)
	}

	for _, o := range options {
		o.AfterParse()
	}
}
