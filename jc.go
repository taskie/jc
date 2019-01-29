package jc

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/yaml.v2"
)

var (
	Version  = "0.1.0-beta"
	Revision = ""
)

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

type Decoder struct {
	Reader io.Reader
	Type   string
}

func NewDecoder(r io.Reader, ty string) *Decoder {
	return &Decoder{
		Reader: r,
		Type:   ty,
	}
}

func (jcDec *Decoder) Decode(data interface{}) error {
	switch strings.ToLower(jcDec.Type) {
	case "json":
		dec := json.NewDecoder(jcDec.Reader)
		err := dec.Decode(data)
		return err
	case "toml":
		_, err := toml.DecodeReader(jcDec.Reader, data)
		return err
	case "yaml":
		dec := yaml.NewDecoder(jcDec.Reader)
		err := dec.Decode(data)
		if err != nil {
			return err
		}
		if v, ok := data.(*interface{}); ok {
			*v = cleanDataFromYaml(data)
		}
		return nil
	case "msgpack":
		dec := msgpack.NewDecoder(jcDec.Reader)
		err := dec.Decode(data)
		return err
	default:
		return fmt.Errorf("invalid --from type: %s", jcDec.Type)
	}
}

func DecodeFile(fpath string, ty string, data interface{}) error {
	f, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer f.Close()
	if ty == "" {
		ty = ExtToTypeMap(filepath.Ext(fpath))
	}
	return NewDecoder(f, ty).Decode(data)
}

type Encoder struct {
	Writer io.Writer
	Type   string
	Indent *string
}

func NewEncoder(w io.Writer, ty string) *Encoder {
	return &Encoder{
		Writer: w,
		Type:   ty,
	}
}

func (jcEnc *Encoder) Encode(data interface{}) error {
	switch strings.ToLower(jcEnc.Type) {
	case "json":
		enc := json.NewEncoder(jcEnc.Writer)
		enc.SetEscapeHTML(false)
		if jcEnc.Indent != nil {
			enc.SetIndent("", *jcEnc.Indent)
		}
		err := enc.Encode(data)
		return err
	case "toml":
		enc := toml.NewEncoder(jcEnc.Writer)
		if jcEnc.Indent != nil {
			enc.Indent = *jcEnc.Indent
		}
		err := enc.Encode(data)
		return err
	case "yaml":
		enc := yaml.NewEncoder(jcEnc.Writer)
		err := enc.Encode(data)
		return err
	case "msgpack":
		enc := msgpack.NewEncoder(jcEnc.Writer)
		err := enc.Encode(data)
		return err
	default:
		return fmt.Errorf("invalid --to type: %s", jcEnc.Type)
	}
}

func EncodeFile(fpath string, ty string, data interface{}) error {
	f, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer f.Close()
	if ty == "" {
		ty = ExtToTypeMap(filepath.Ext(fpath))
	}
	return NewEncoder(f, ty).Encode(data)
}

var extToTypeMap = map[string]string{
	".json":        "json",
	".toml":        "toml",
	".yml":         "yaml",
	".yaml":        "yaml",
	".mp":          "msgpack",
	".msgpack":     "msgpack",
	".messagepack": "msgpack",
}

func ExtToTypeMap(ext string) string {
	ty := extToTypeMap[ext]
	if ty != "" {
		return ty
	}
	return "json"
}

type Jc struct {
	FromType string
	ToType   string
	Indent   *string
}

func (jc *Jc) Decode(r io.Reader, data interface{}) error {
	return NewDecoder(r, jc.FromType).Decode(data)
}

func (jc *Jc) Encode(w io.Writer, data interface{}) error {
	return NewEncoder(w, jc.ToType).Encode(data)
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
