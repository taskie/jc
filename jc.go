package jc

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/yaml.v2"
	"io"
	"strings"
)

var (
	Version  = "0.1.0-beta"
	Revision = ""
)

type Jc struct {
	FromType string
	ToType   string
	Indent   *string
}

func cleanDataFromYaml(data interface{}) interface{} {
	switch oldData := data.(type) {
	case *interface{}:
		return cleanDataFromYaml(*oldData)
	case map[interface{}]interface{}:
		newData := make(map[string]interface{})
		for k, v := range oldData {
			s := fmt.Sprintf("%v", k)
			newData[s] = cleanDataFromYaml(v)
		}
		return newData
	case []interface{}:
		newData := make([]interface{}, len(oldData))
		for i, v := range oldData {
			newData[i] = cleanDataFromYaml(v)
		}
		return newData
	default:
		return data
	}
}

func (jc *Jc) Decode(r io.Reader, data interface{}) error {
	switch strings.ToLower(jc.FromType) {
	case "json":
		dec := json.NewDecoder(r)
		err := dec.Decode(data)
		return err
	case "toml":
		_, err := toml.DecodeReader(r, data)
		return err
	case "yaml":
		dec := yaml.NewDecoder(r)
		err := dec.Decode(data)
		if err != nil {
			return err
		}
		if v, ok := data.(*interface{}); ok {
			*v = cleanDataFromYaml(data)
		}
		return nil
	case "msgpack":
		dec := msgpack.NewDecoder(r)
		err := dec.Decode(data)
		return err
	default:
		return fmt.Errorf("invalid --from type: %s", jc.FromType)
	}
}

func (jc *Jc) Encode(w io.Writer, data interface{}) error {
	switch strings.ToLower(jc.ToType) {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		if jc.Indent != nil {
			enc.SetIndent("", *jc.Indent)
		}
		err := enc.Encode(data)
		return err
	case "toml":
		enc := toml.NewEncoder(w)
		if jc.Indent != nil {
			enc.Indent = *jc.Indent
		}
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
		return fmt.Errorf("invalid --to type: %s", jc.ToType)
	}
}

func (jc *Jc) Run(r io.Reader, w io.Writer) error {
	var data interface{}
	err := jc.Decode(r, &data)
	if err != nil {
		return err
	}
	err = jc.Encode(w, data)
	return err
}
