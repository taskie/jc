package jc

import (
	"io"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/taskie/fwv"
	"github.com/taskie/jc"
	"github.com/taskie/ose"
	"github.com/taskie/ose/coli"
	"go.uber.org/zap"
)

const CommandName = "jc"

var Command *cobra.Command

func init() {
	Command = NewCommand(coli.NewColiInThisWorld())
}

func Main() {
	Command.Execute()
}

func NewCommand(cl *coli.Coli) *cobra.Command {
	cmd := &cobra.Command{
		Use:  CommandName,
		Args: cobra.RangeArgs(0, 2),
		Run:  cl.WrapRun(run),
	}
	cl.Prepare(cmd)

	flg := cmd.Flags()
	flg.StringP("from-type", "f", "", "convert from [json|toml|yaml|msgpack|dotenv]")
	flg.StringP("to-type", "t", "", "convert to [json|toml|yaml|msgpack|dotenv]")
	flg.StringP("indent", "I", "", "indentation of output")

	cl.BindFlags(flg, []string{"from-type", "to-type", "indent"})
	return cmd
}

type Config struct {
	FromType, ToType, Indent, LogLevel string
}

func run(cl *coli.Coli, cmd *cobra.Command, args []string) {
	v := cl.Viper()
	log := zap.L()
	if v.GetBool("version") {
		cmd.Println(fwv.Version)
		return
	}
	var config Config
	err := v.Unmarshal(&config)
	if err != nil {
		log.Fatal("can't unmarshal config", zap.Error(err))
	}

	input := ""
	output := ""
	switch len(args) {
	case 0:
		break
	case 1:
		input = args[0]
	case 2:
		input = args[0]
		output = args[1]
	default:
		log.Fatal("invalid arguments", zap.Strings("arguments", args[2:]))
	}

	fromType := config.FromType
	if fromType == "" {
		fromType = jc.ExtToType(filepath.Ext(input))
	}
	toType := config.ToType
	if toType == "" {
		toType = jc.ExtToType(filepath.Ext(output))
	}

	opener := ose.NewOpenerInThisWorld()
	r, err := opener.Open(input)
	if err != nil {
		log.Fatal("can't open", zap.Error(err))
	}
	defer r.Close()
	_, err = opener.CreateTempFile("", CommandName, output, func(f io.WriteCloser) (bool, error) {
		jc := jc.Converter{
			FromType: fromType,
			ToType:   toType,
			Indent:   &config.Indent,
		}
		err = jc.Convert(f, r)
		if err != nil {
			log.Fatal("can't convert", zap.Error(err))
		}
		return true, nil
	})
	if err != nil {
		log.Fatal("can't create file", zap.Error(err))
	}
}
