package proto

import (
	"bytes"
	"testing"
)

func TestGenerateCodeBytesArray(t *testing.T) {
	b, err := GenerateCodeBytesArray("600000")
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(b[:], []byte{'6', '0', '0', '0', '0', '0'}) {
		t.Error("GenerateCodeBytesArray error")
		return
	}
}

func TestNewTDXCodec(t *testing.T) {
	cc, err := NewTDXCodec()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("\n%s", cc.dumpBook())
}
