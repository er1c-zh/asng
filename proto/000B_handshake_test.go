package proto

import (
	"encoding/hex"
	"testing"
)

func TestHandshake(t *testing.T) {
	handshake := &Handshake{}
	handshake.Req = &HandshakeReq{}
	handshake.Req.R0 = 0x0
	handshake.Req.LoginType = [50]byte{'g', 'u', 'e', 's', 't'}
	handshake.Req.Flag0 = 0x7
	handshake.Req.Flag1 = 0x1
	handshake.Req.Flag2 = 0x0104
	handshake.Req.Flag3 = 0x0
	handshake.Req.Version = [4]byte{0xB8, 0x1E, 0xF5, 0x40}
	handshake.Req.MinVersion = [4]byte{0xAE, 0x47, 0xC1, 0x40}
	handshake.Req.Flag4 = 0x1
	handshake.Req.Flag5 = [2]uint8{0x1, 0x1}
	handshake.Req.SomeKey = [12]byte{'0', '0', '0', 'C', '2', '9', 'F', '0', '1', '9', 'C', '2'}
	handshake.Req.Flag6 = 0x000056C2

	d, err := handshake.MarshalReqBody(nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("\n%s", hex.Dump(d))

	ec, err := NewTDXCodec()
	if err != nil {
		t.Fatal(err)
		return
	}
	nd, err := ec.Encode(d)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("\n%s", hex.Dump(nd))

	nd, err = ec.Encode(nd)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("\n%s", hex.Dump(nd))

	nd, err = ec.Encode(nd)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("\n%s", hex.Dump(nd))
}
