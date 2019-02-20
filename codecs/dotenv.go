package codecs

import (
	"fmt"
	"io"
	"reflect"

	"github.com/joho/godotenv"
)

var (
	DotenvCodec = &Codec{
		Type: "dotenv",
		Exts: []string{".env"},
		DecoderBuilder: func(r io.Reader, opts interface{}) (Decoder, error) {
			return NewDotenvDecoder(r), nil
		},
		EncoderBuilder: func(w io.Writer, opts interface{}) (Encoder, error) {
			return NewDotenvEncoder(w), nil
		},
	}
)

type DotenvDecoder struct {
	reader io.Reader
}

func NewDotenvDecoder(r io.Reader) *DotenvDecoder {
	return &DotenvDecoder{
		reader: r,
	}
}

func (d *DotenvDecoder) Decode(data interface{}) error {
	m, err := godotenv.Parse(d.reader)
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
}

type DotenvEncoder struct {
	writer io.Writer
}

func NewDotenvEncoder(w io.Writer) *DotenvEncoder {
	return &DotenvEncoder{
		writer: w,
	}
}

func (e *DotenvEncoder) Encode(data interface{}) error {
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
	_, err = io.WriteString(e.writer, s)
	return err
}
