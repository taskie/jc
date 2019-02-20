package jc

import (
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/taskie/jc"
	"github.com/taskie/osplus"
)

var (
	cfgFile, fromType, toType, indent string
	noColor, verbose, version         bool
)

func init() {
	cobra.OnInitialize(initConfig)
	Command.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $XDG_CONFIG_HOME/jc/jc.yml)")
	Command.Flags().StringVarP(&fromType, "fromType", "f", "", "convert from [json|toml|yaml|msgpack|dotenv]")
	Command.Flags().StringVarP(&toType, "toType", "t", "", "convert to [json|toml|yaml|msgpack|dotenv]")
	Command.Flags().StringVarP(&indent, "indent", "I", "", "indentation of output")
	Command.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	Command.Flags().BoolVarP(&version, "version", "V", false, "show Version")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		conf, err := osplus.GetXdgConfigHome()
		if err != nil {
			panic(err)
		}
		viper.AddConfigPath(filepath.Join(conf, "jc"))
		viper.SetConfigName("jc")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func Main() {
	Command.Execute()
}

var Command = &cobra.Command{
	Use: "jc",
	Run: func(cmd *cobra.Command, args []string) {
		err := run(cmd, args)
		if err != nil {
			logrus.Fatal(err)
		}
	},
}

func run(cmd *cobra.Command, args []string) error {
	if version {
		fmt.Println(jc.Version)
		return nil
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
		return fmt.Errorf("invalid arguments: %v", args[2:])
	}

	if fromType == "" {
		fromType = jc.ExtToType(filepath.Ext(input))
	}
	if toType == "" {
		toType = jc.ExtToType(filepath.Ext(output))
	}

	opener := osplus.NewOpener()
	r, err := opener.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()
	w, commit, err := opener.CreateTempFileWithDestination(output, "", "jc-")
	if err != nil {
		return err
	}
	defer w.Close()

	jc := jc.Converter{
		FromType: fromType,
		ToType:   toType,
		Indent:   &indent,
	}
	err = jc.Convert(w, r)
	if err != nil {
		return err
	}
	commit(true)
	return nil
}
