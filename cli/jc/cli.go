package jc

import (
	"fmt"
	"path/filepath"

	"github.com/iancoleman/strcase"
	"github.com/k0kubun/pp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/taskie/jc"
	"github.com/taskie/osplus"
)

type Config struct {
	FromType, ToType, Indent, LogLevel string
}

var configFile string
var config Config
var (
	verbose, debug, version bool
)

const CommandName = "jc"

func init() {
	Command.PersistentFlags().StringVarP(&configFile, "config", "c", "", `config file (default "jc.yml")`)
	Command.Flags().StringP("from-type", "f", "", "convert from [json|toml|yaml|msgpack|dotenv]")
	Command.Flags().StringP("to-type", "t", "", "convert to [json|toml|yaml|msgpack|dotenv]")
	Command.Flags().StringP("indent", "I", "", "indentation of output")
	Command.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	Command.Flags().BoolVarP(&debug, "debug", "g", false, "debug output")
	Command.Flags().BoolVarP(&version, "version", "V", false, "show Version")

	for _, s := range []string{"from-type", "to-type", "indent"} {
		envKey := strcase.ToSnake(s)
		structKey := strcase.ToCamel(s)
		viper.BindPFlag(envKey, Command.Flags().Lookup(s))
		viper.RegisterAlias(structKey, envKey)
	}

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if debug {
		log.SetLevel(log.DebugLevel)
	} else if verbose {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName(CommandName)
		conf, err := osplus.GetXdgConfigHome()
		if err != nil {
			log.Info(err)
		} else {
			viper.AddConfigPath(filepath.Join(conf, CommandName))
		}
		viper.AddConfigPath(".")
	}
	viper.SetEnvPrefix(CommandName)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Debug(err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Warn(err)
	}
}

func Main() {
	Command.Execute()
}

var Command = &cobra.Command{
	Use:  CommandName,
	Args: cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		err := run(cmd, args)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func run(cmd *cobra.Command, args []string) error {
	if version {
		fmt.Println(jc.Version)
		return nil
	}
	if config.LogLevel != "" {
		lv, err := log.ParseLevel(config.LogLevel)
		if err != nil {
			log.Warn(err)
		} else {
			log.SetLevel(lv)
		}
	}
	if debug {
		if viper.ConfigFileUsed() != "" {
			log.Debugf("Using config file: %s", viper.ConfigFileUsed())
		}
		log.Debug(pp.Sprint(config))
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

	fromType := config.FromType
	if fromType == "" {
		fromType = jc.ExtToType(filepath.Ext(input))
	}
	toType := config.ToType
	if toType == "" {
		toType = jc.ExtToType(filepath.Ext(output))
	}

	opener := osplus.NewOpener()
	r, err := opener.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()
	w, commit, err := opener.CreateTempFileWithDestination(output, "", CommandName+"-")
	if err != nil {
		return err
	}
	defer w.Close()

	jc := jc.Converter{
		FromType: fromType,
		ToType:   toType,
		Indent:   &config.Indent,
	}
	err = jc.Convert(w, r)
	if err != nil {
		return err
	}
	commit(true)
	return nil
}
