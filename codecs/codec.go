package codecs

import (
	"io"
	"strings"
)

type Decoder interface {
	Decode(data interface{}) error
}

type Encoder interface {
	Encode(data interface{}) error
}

type Codec struct {
	Type           string
	Aliases        []string
	Exts           []string
	DecoderBuilder func(r io.Reader, opts interface{}) (Decoder, error)
	EncoderBuilder func(w io.Writer, opts interface{}) (Encoder, error)
}

type CodecResolver struct {
	codecMap       map[string]*Codec
	aliasToTypeMap map[string]string
	extToTypeMap   map[string]string
}

func NewCodecResolver() *CodecResolver {
	return &CodecResolver{
		codecMap:       make(map[string]*Codec),
		aliasToTypeMap: make(map[string]string),
		extToTypeMap:   make(map[string]string),
	}
}

func (r *CodecResolver) Register(ty string, codec *Codec) {
	r.codecMap[ty] = codec
}

func (r *CodecResolver) RegisterAll(codec *Codec) {
	r.Register(codec.Type, codec)
	if codec.Exts != nil {
		for _, ext := range codec.Exts {
			r.extToTypeMap[ext] = codec.Type
		}
	}
	if codec.Aliases != nil {
		for _, alias := range codec.Aliases {
			r.aliasToTypeMap[alias] = codec.Type
		}
	}
}

func (r *CodecResolver) RegisterAsDefault(codec *Codec) {
	r.Register("", codec)
}

func (r *CodecResolver) Resolve(key string) *Codec {
	key = strings.ToLower(key)
	aliasKey := r.aliasToTypeMap[key]
	if aliasKey != "" {
		key = aliasKey
	}
	return r.codecMap[key]
}

func (r *CodecResolver) ResolveWithExt(ext string) *Codec {
	return r.codecMap[r.ExtToType(ext)]
}

func (r *CodecResolver) ExtToType(ext string) string {
	return r.extToTypeMap[ext]
}
