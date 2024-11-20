package proto

import (
	"bytes"
	"context"
	"encoding/binary"
)

func (c *Client) Subscribe(market uint8, code string) error {
	subject := &Subscribe{}

	subject.SetDebug(c.ctx)

	subject.Req = &SubscribeReq{
		R0:     0,
		Market: market,
		Code:   [6]byte{code[0], code[1], code[2], code[3], code[4], code[5]},
		R1:     [4]uint8{0x00, 0x00, 0x00, 0x39}, // 0x00000039
	}
	copy(subject.Req.Code[:], code)

	err := do(c, c.dataConn, subject)
	if err != nil {
		return err
	}
	return nil
}

type Subscribe struct {
	BlankCodec
	Req *SubscribeReq
}

type SubscribeReq struct {
	R0     uint8
	Market uint8
	Code   [6]byte
	R1     [4]uint8
}

func (Subscribe) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.Method = 0x0537
	return nil
}

func (obj *Subscribe) MarshalReqBody(ctx context.Context) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := binary.Write(buf, binary.LittleEndian, obj.Req)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (obj *Subscribe) UnmarshalResp(ctx context.Context, data []byte) error {
	return nil
}
