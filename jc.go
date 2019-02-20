package jc

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/taskie/jc/codecs"
)

var (
	Version = "0.1.0-beta"
)

var defaultResolver *codecs.CodecResolver

func init() {
	defaultResolver = codecs.NewCodecResolver()
	defaultResolver.RegisterAsDefault(codecs.JsonCodec)
	defaultResolver.RegisterAll(codecs.JsonCodec)
	defaultResolver.RegisterAll(codecs.MsgpackCodec)
	defaultResolver.RegisterAll(codecs.YamlCodec)
	defaultResolver.RegisterAll(codecs.TomlCodec)
	defaultResolver.RegisterAll(codecs.DotenvCodec)
}

func ExtToType(ext string) string {
	return defaultResolver.ExtToType(ext)
}

type Decoder struct {
	Resolver *codecs.CodecResolver
	Reader   io.Reader
	Type     string
}

func NewDecoder(r io.Reader, ty string) *Decoder {
	return &Decoder{
		Resolver: defaultResolver,
		Reader:   r,
		Type:     ty,
	}
}

func (d *Decoder) Decode(data interface{}) error {
	codec := d.Resolver.Resolve(d.Type)
	if codec == nil {
		return fmt.Errorf("invalid input type: %s", d.Type)
	}
	dec, err := codec.DecoderBuilder(d.Reader, nil)
	if err != nil {
		return err
	}
	return dec.Decode(data)
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
	Resolver *codecs.CodecResolver
	Writer   io.Writer
	Type     string
	Indent   *string
}

func NewEncoder(w io.Writer, ty string) *Encoder {
	return &Encoder{
		Resolver: defaultResolver,
		Writer:   w,
		Type:     ty,
	}
}

func (e *Encoder) Encode(data interface{}) error {
	codec := e.Resolver.Resolve(e.Type)
	if codec == nil {
		return fmt.Errorf("invalid output type: %s", e.Type)
	}
	opts := map[string]interface{}{
		"escapeHTML": false,
	}
	if e.Indent != nil {
		opts["indent"] = *e.Indent
	}
	enc, err := codec.EncoderBuilder(e.Writer, opts)
	if err != nil {
		return err
	}
	return enc.Encode(data)
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

type Converter struct {
	Resolver *codecs.CodecResolver
	FromType string
	ToType   string
	Indent   *string
}

func (c *Converter) Decode(r io.Reader, data interface{}) error {
	d := NewDecoder(r, c.FromType)
	if c.Resolver != nil {
		d.Resolver = c.Resolver
	}
	return d.Decode(data)
}

func (c *Converter) Encode(w io.Writer, data interface{}) error {
	e := NewEncoder(w, c.ToType)
	e.Indent = c.Indent
	if c.Resolver != nil {
		e.Resolver = c.Resolver
	}
	return e.Encode(data)
}

func (c *Converter) Convert(dst io.Writer, src io.Reader) error {
	var data interface{}
	err := c.Decode(src, &data)
	if err != nil {
		return err
	}
	err = c.Encode(dst, data)
	return err
}
