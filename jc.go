package jc

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strings"
	"github.com/vmihailenco/msgpack"
	"github.com/jessevdk/go-flags"
)

var (
	version  string
	revision string
)

type Jc struct {
	FromType string
	ToType   string
	Options  map[string]interface{}
}

func Decode(r io.Reader, fromType string, data interface{}) error {
	switch strings.ToLower(fromType) {
	case "json":
		dec := json.NewDecoder(r)
		err := dec.Decode(data)
		return err
	case "toml":
		_, err := toml.DecodeReader(r, data)
		return err
	//case "yaml":
	//	dec := yaml.NewDecoder(r)
	//	err := dec.Decode(data)
	//	return err
	case "msgpack":
		dec := msgpack.NewDecoder(r)
		err := dec.Decode(data)
		return err
	default:
		return fmt.Errorf("invalid --from type: %s", fromType)
	}
}

func Encode(w io.Writer, toType string, data interface{}) error {
	switch strings.ToLower(toType) {
	case "json":
		enc := json.NewEncoder(w)
		err := enc.Encode(data)
		return err
	case "toml":
		enc := toml.NewEncoder(w)
		err := enc.Encode(data)
		return err
	case "yaml":
		enc := yaml.NewEncoder(w)
		err := enc.Encode(data)
		return err
	case "msgpack":
		enc := msgpack.NewEncoder(w)
		err := enc.Encode(data)
		return err
	default:
		return fmt.Errorf("invalid --to type: %s", toType)
	}
}

func (jc *Jc) Run(r io.Reader, w io.Writer) error {
	var data interface{}
	err := Decode(r, jc.FromType, &data)
	if err != nil {
		return err
	}
	err = Encode(w, jc.ToType, data)
	return err
}

type Options struct {
	FromType string `short:"f" long:"from" default:"json" description:"convert from [json|toml|msgpack]"`
	ToType   string `short:"t" long:"to" default:"json" description:"convert to [json|toml|yaml|msgpack]"`
	// Indent   string `short:"I" long:"indent" description:"indentation of output"`
	NoColor bool `long:"no-color" env:"NO_COLOR" description:"NOT colorize output"`
	Verbose bool `short:"v" long:"verbose" description:"show verbose output"`
	Version bool `short:"V" long:"version" description:"show version"`
}

func Main(args []string) {
	var opts Options
	args, err := flags.ParseArgs(&opts, args)
	if opts.Version {
		if opts.Verbose {
			fmt.Println("Version: ", version)
			fmt.Println("Revision: ", revision)
		} else {
			fmt.Println(version)
		}
		os.Exit(0)
	}

	jc := Jc{
		FromType: opts.FromType,
		ToType:   opts.ToType,
	}
	err = jc.Run(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
