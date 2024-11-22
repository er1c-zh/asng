package proto

import (
	"asng/models"
	"bytes"
	"context"
	"encoding/binary"
)

func (c *Client) Subscribe(market models.MarketType, code string) (*SubscribeResp, error) {
	subject := &Subscribe{}

	subject.SetDebug(c.ctx)

	subject.Req = &SubscribeReq{
		Market: market,
		Code:   [6]byte{code[0], code[1], code[2], code[3], code[4], code[5]},
		R1:     [4]uint8{0x00, 0x00, 0x00, uint8(c.handShakeSeed)},
	}

	err := do(c, c.dataConn, subject)
	if err != nil {
		return nil, err
	}
	return subject.Resp, nil
}

type Subscribe struct {
	BlankCodec
	Req  *SubscribeReq
	Resp *SubscribeResp
}

type SubscribeReq struct {
	Market models.MarketType
	Code   [6]byte
	R1     [4]uint8
}
type SubscribeResp struct {
	Count uint16
	List  []QuoteFrame
}

type QuoteFrame struct {
	Price    int64
	AvgPrice int64
	Volume   int64
}

func (Subscribe) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.Type0 = 0x0100
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
	var err error
	obj.Resp = &SubscribeResp{}
	cursor := 0
	obj.Resp.Count, err = ReadInt(data, &cursor, obj.Resp.Count)
	if err != nil {
		return err
	}
	if obj.Resp.Count <= 0 {
		return nil
	}
	base0, base1 := int64(0), int64(0)
	cursor += 2
	for i := 0; i < int(obj.Resp.Count); i += 1 {
		item := QuoteFrame{}
		item.Price, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		item.Price = item.Price + base0
		item.AvgPrice, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		item.AvgPrice = item.AvgPrice + base1
		item.Volume, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		if i == 0 {
			base0 = item.Price
			base1 = item.AvgPrice
		}
		obj.Resp.List = append(obj.Resp.List, item)
	}

	return nil
}
