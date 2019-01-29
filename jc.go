package jc

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
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
		} else {
			panic("data must be pointer")
		}
		return nil
	case "msgpack":
		dec := msgpack.NewDecoder(jcDec.Reader)
		err := dec.Decode(data)
		return err
	case "dotenv":
		m, err := godotenv.Parse(jcDec.Reader)
		if err != nil {
			return err
		}
		if v, ok := data.(*interface{}); ok {
			*v = m
		} else if v, ok := data.(*map[string]string); ok {
			*v = m
		} else {
			panic("data must be pointer of string map or interface{}")
		}
		return nil
	default:
		return fmt.Errorf("invalid input type: %s", jcDec.Type)
	}
}

func DecodeFile(fpath string, ty string, data interface{}) error {
	f, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer f.Close()
	if ty == "" {
		ty = ExtToType(filepath.Ext(fpath))
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
	case "dotenv":
		m := make(map[string]string)
		rv := reflect.ValueOf(data)
		if rv.Kind() != reflect.Map {
			return fmt.Errorf("value is not map")
		}
		ks := rv.MapKeys()
		for _, k := range ks {
			var sk string
			var ok bool
			if sk, ok = k.Interface().(string); !ok {
				return fmt.Errorf("key is not string")
			}
			v := rv.MapIndex(k)
			var sv string
			if sv, ok = v.Interface().(string); !ok {
				return fmt.Errorf("value is not string")
			}
			m[sk] = sv
		}
		s, err := godotenv.Marshal(m)
		if err != nil {
			return err
		}
		_, err = io.WriteString(jcEnc.Writer, s)
		return err
	default:
		return fmt.Errorf("invalid output type: %s", jcEnc.Type)
	}
}

func EncodeFile(fpath string, ty string, data interface{}) error {
	f, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer f.Close()
	if ty == "" {
		ty = ExtToType(filepath.Ext(fpath))
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
	".env":         "dotenv",
}

func ExtToType(ext string) string {
	ty := extToTypeMap[strings.ToLower(ext)]
	if ty != "" {
		return ty
	}
	return "json"
}

type Converter struct {
	FromType string
	ToType   string
	Indent   *string
}

func (conv *Converter) Decode(r io.Reader, data interface{}) error {
	return NewDecoder(r, conv.FromType).Decode(data)
}

func (conv *Converter) Encode(w io.Writer, data interface{}) error {
	return NewEncoder(w, conv.ToType).Encode(data)
}

func (conv *Converter) Convert(dst io.Writer, src io.Reader) error {
	var data interface{}
	err := conv.Decode(src, &data)
	if err != nil {
		return err
	}
	err = conv.Encode(dst, data)
	return err
}
