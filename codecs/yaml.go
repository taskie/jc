package codecs

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

var (
	YamlCodec = &Codec{
		Type: "yaml",
		Exts: []string{".yml", ".yaml"},
		DecoderBuilder: func(r io.Reader, opts interface{}) (Decoder, error) {
			return NewYamlDecoder(r), nil
		},
		EncoderBuilder: func(w io.Writer, opts interface{}) (Encoder, error) {
			return NewYamlEncoder(w), nil
		},
	}
)

type YamlDecoder struct {
	underlyingDecoder *yaml.Decoder
}

func NewYamlDecoder(r io.Reader) *YamlDecoder {
	return &YamlDecoder{
		underlyingDecoder: yaml.NewDecoder(r),
	}
}

func (d *YamlDecoder) Decode(data interface{}) error {
	err := d.underlyingDecoder.Decode(data)
	if err != nil {
		return err
	}
	if v, ok := data.(*interface{}); ok {
		*v = cleanDataFromYaml(data)
	} else {
		return fmt.Errorf("data must be pointer")
	}
	return nil
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
	case map[string]interface{}:
		newData := make(map[string]interface{})
		for k, v := range oldData {
			newData[k] = cleanDataFromYaml(v)
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

type YamlEncoder = yaml.Encoder

func NewYamlEncoder(w io.Writer) *YamlEncoder {
	return yaml.NewEncoder(w)
}
