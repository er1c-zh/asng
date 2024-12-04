package proto

import (
	"bytes"
	"fmt"
	"testing"
	"time"
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

func TestFormatDate(t *testing.T) {
	fmt.Printf("%s\n", time.Date(2024, 9, 27, 0, 0, 0, 0, time.Local).AddDate(0, 0, -55839).Format("2006-01-02"))
	fmt.Printf("%x\n", int64(time.Date(2021, 2, 1, 0, 0, 0, 0, time.Local).Sub(time.Date(1871, 11, 10, 0, 0, 0, 0, time.Local)).Hours()/24))
}
