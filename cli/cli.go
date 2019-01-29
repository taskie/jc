package cli

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/taskie/osplus"

	"github.com/jessevdk/go-flags"
	"github.com/taskie/jc"
)

type Options struct {
	FromType string  `short:"f" long:"from" description:"convert from [json|toml|yaml|msgpack|dotenv]"`
	ToType   string  `short:"t" long:"to" description:"convert to [json|toml|yaml|msgpack|dotenv]"`
	Indent   *string `short:"I" long:"indent" description:"indentation of output"`
	NoColor  bool    `long:"noColor" env:"NO_COLOR" description:"NOT colorize output"`
	Verbose  bool    `short:"v" long:"verbose" description:"show verbose output"`
	Version  bool    `short:"V" long:"version" description:"show version"`
}

func mainImpl() error {
	var opts Options
	args, err := flags.ParseArgs(&opts, os.Args)
	if opts.Version {
		if opts.Verbose {
			fmt.Println("Version: ", jc.Version)
			fmt.Println("Revision: ", jc.Revision)
		} else {
			fmt.Println(jc.Version)
		}
		os.Exit(0)
	}

	var tmpFile *os.File
	var outFilePath string
	var w io.Writer
	var r io.Reader
	fromType := opts.FromType
	toType := opts.ToType

	if len(args) > 3 {
		return fmt.Errorf("invalid arguments")
	}
	if len(args) > 2 && args[2] != "-" {
		outFilePath = args[2]
		tmpFile, err = ioutil.TempFile("", "jc-")
		if err != nil {
			return err
		}
		defer func() {
			if tmpFile != nil {
				tmpFile.Close()
				log.Print(tmpFile.Name())
				os.Remove(tmpFile.Name())
			}
		}()
		w = tmpFile
		toType = jc.ExtToType(filepath.Ext(outFilePath))
	} else {
		w = os.Stdout
	}
	if len(args) > 1 && args[1] != "-" {
		inFilePath := args[1]
		inFile, err := os.Open(inFilePath)
		if err != nil {
			return err
		}
		defer inFile.Close()
		r = inFile
		fromType = jc.ExtToType(filepath.Ext(inFilePath))
	} else {
		r = os.Stdout
	}

	if fromType == "" {
		fromType = "json"
	}
	if toType == "" {
		toType = "json"
	}

	jc := jc.Converter{
		FromType: fromType,
		ToType:   toType,
		Indent:   opts.Indent,
	}
	err = jc.Convert(w, r)
	if err != nil {
		return err
	}
	if tmpFile != nil {
		tmpName := tmpFile.Name()
		err = tmpFile.Close()
		if err != nil {
			return err
		}
		err = osplus.MoveFile(tmpName, outFilePath, &osplus.MoveOptions{NoOverwrite: true})
		if err != nil {
			return err
		}
		tmpFile = nil
	}
	return nil
}

func Main() {
	err := mainImpl()
	if err != nil {
		log.Fatal(err)
	}
}
