package codecs

import (
	"encoding/json"
	"io"
)

var (
	JsonCodec = &Codec{
		Type: "json",
		Exts: []string{".json"},
		DecoderBuilder: func(r io.Reader, opts interface{}) (Decoder, error) {
			return NewJsonDecoder(r), nil
		},
		EncoderBuilder: func(w io.Writer, opts interface{}) (Encoder, error) {
			enc := NewJsonEncoder(w)
			if v, ok := opts.(map[string]interface{})["indent"].(string); ok {
				enc.SetIndent("", v)
			}
			if v, ok := opts.(map[string]interface{})["escapeHTML"].(bool); ok {
				enc.SetEscapeHTML(v)
			}
			return enc, nil
		},
	}
)

type JsonDecoder = json.Decoder

func NewJsonDecoder(r io.Reader) *JsonDecoder {
	return json.NewDecoder(r)
}

type JsonEncoder = json.Encoder

func NewJsonEncoder(w io.Writer) *JsonEncoder {
	return json.NewEncoder(w)
}
