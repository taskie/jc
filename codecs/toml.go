package codecs

import (
	"io"

	"github.com/BurntSushi/toml"
)

var (
	TomlCodec = &Codec{
		Type: "toml",
		Exts: []string{".toml"},
		DecoderBuilder: func(r io.Reader, opts interface{}) (Decoder, error) {
			return NewTomlDecoder(r), nil
		},
		EncoderBuilder: func(w io.Writer, opts interface{}) (Encoder, error) {
			enc := NewTomlEncoder(w)
			if v, ok := opts.(map[string]interface{})["indent"].(string); ok {
				enc.Indent = v
			}
			return enc, nil
		},
	}
)

type TomlDecoder struct {
	reader io.Reader
}

func NewTomlDecoder(r io.Reader) *TomlDecoder {
	return &TomlDecoder{
		reader: r,
	}
}

func (d *TomlDecoder) Decode(data interface{}) error {
	_, err := toml.DecodeReader(d.reader, data)
	return err
}

type TomlEncoder = toml.Encoder

func NewTomlEncoder(w io.Writer) *TomlEncoder {
	return toml.NewEncoder(w)
}
