package codecs

import (
	"io"

	"github.com/vmihailenco/msgpack"
)

var (
	MsgpackCodec = &Codec{
		Type:    "msgpack",
		Aliases: []string{"msgpack", "messagepack"},
		Exts:    []string{".msgpack", ".messagepack", ".mpack", ".mp"},
		DecoderBuilder: func(r io.Reader, opts interface{}) (Decoder, error) {
			return NewMsgpackDecoder(r), nil
		},
		EncoderBuilder: func(w io.Writer, opts interface{}) (Encoder, error) {
			return NewMsgpackEncoder(w), nil
		},
	}
)

type MsgpackDecoder = msgpack.Decoder

func NewMsgpackDecoder(r io.Reader) *MsgpackDecoder {
	return msgpack.NewDecoder(r)
}

type MsgpackEncoder = msgpack.Encoder

func NewMsgpackEncoder(w io.Writer) *MsgpackEncoder {
	return msgpack.NewEncoder(w)
}
