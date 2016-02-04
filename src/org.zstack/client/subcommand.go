package client

import (
	"flag"
	"fmt"
	"os"
)

type (
	SubCommand interface {
		Name() string
		Flags(*flag.FlagSet)
		CheckFlags() error
		Run() int
	}

	_subCommandWrapper struct {
		cmd     SubCommand
		flagSet *flag.FlagSet
	}
)

var (
	_subCommands   []SubCommand = make([]SubCommand, 0)
	_subCommandMap map[string]*_subCommandWrapper
)

func RegisterSubcommand(c SubCommand) {
	_subCommands = append(_subCommands, c)
}

func ParseSubCommands() {
	_subCommandMap = make(map[string]*_subCommandWrapper, len(_subCommands))
	for _, sc := range _subCommands {
		name := sc.Name()
		w := &_subCommandWrapper{
			cmd:     sc,
			flagSet: flag.NewFlagSet(name, flag.ExitOnError),
		}

		sc.Flags(w.flagSet)
		_subCommandMap[name] = w
	}

	oldUsage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n\n")
		oldUsage()
		for name, sc := range _subCommandMap {
			fmt.Fprintf(os.Stderr, "\n# %s %s\n", os.Args[0], name)
			sc.flagSet.PrintDefaults()
			fmt.Fprintf(os.Stderr, "\n")
		}
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

}

func RunSubCommand() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()

	cmdName := flag.Arg(0)

	if sc, ok := _subCommandMap[cmdName]; ok {
		sc.flagSet.Parse(flag.Args()[1:])

		if err := sc.cmd.CheckFlags(); err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			flag.Usage()
			os.Exit(1)
		}

		ret := sc.cmd.Run()
		os.Exit(ret)
	} else {
		fmt.Fprintf(os.Stderr, "error: command[%s] not found\n", cmdName)
		flag.Usage()
		os.Exit(1)
	}
}
