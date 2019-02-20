package codecs

import (
	"bytes"
	"testing"
)

func TestMsgpackCodec(t *testing.T) {
	var m interface{} = map[string]interface{}{
		"x": 42,
		"y": 3.14,
		"z": map[string]interface{}{
			"a": "hello",
			"b": true,
			"c": nil,
		},
	}
	for i := 0; i < 2; i++ {
		bw := &bytes.Buffer{}
		var enc Encoder = NewMsgpackEncoder(bw)
		err := enc.Encode(m)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%v", bw.Bytes())
		br := bytes.NewBuffer(bw.Bytes())
		var dec Decoder = NewMsgpackDecoder(br)
		var n interface{}
		err = dec.Decode(&n)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%+v", n)
		m = n
	}
	// t.Fail()
}
