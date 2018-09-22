package cli

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/taskie/jc"
	"os"
)

type Options struct {
	FromType string  `short:"f" long:"from" default:"json" description:"convert from [json|toml|msgpack]"`
	ToType   string  `short:"t" long:"to" default:"json" description:"convert to [json|toml|yaml|msgpack]"`
	Indent   *string `short:"I" long:"indent" description:"indentation of output"`
	NoColor  bool    `long:"noColor" env:"NO_COLOR" description:"NOT colorize output"`
	Verbose  bool    `short:"v" long:"verbose" description:"show verbose output"`
	Version  bool    `short:"V" long:"version" description:"show version"`
}

func Main() {
	var opts Options
	_, err := flags.ParseArgs(&opts, os.Args)
	if opts.Version {
		if opts.Verbose {
			fmt.Println("Version: ", jc.Version)
			fmt.Println("Revision: ", jc.Revision)
		} else {
			fmt.Println(jc.Version)
		}
		os.Exit(0)
	}

	jc := jc.Jc{
		FromType: opts.FromType,
		ToType:   opts.ToType,
		Indent:   opts.Indent,
	}
	err = jc.Run(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
