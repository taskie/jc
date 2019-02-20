package codecs

import (
	"bytes"
	"testing"
)

func TestDotenvCodec(t *testing.T) {
	var m interface{} = map[string]interface{}{
		"x": "hello",
		"y": "world",
		"z": "!",
	}
	for i := 0; i < 2; i++ {
		bw := &bytes.Buffer{}
		var enc Encoder = NewDotenvEncoder(bw)
		err := enc.Encode(m)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(bw)
		br := bytes.NewBuffer(bw.Bytes())
		var dec Decoder = NewDotenvDecoder(br)
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
