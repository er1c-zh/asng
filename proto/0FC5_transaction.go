package proto

import (
	"bytes"
	"context"
	"encoding/binary"
)

type Transaction struct {
	BlankCodec
	Req        *TransactionReq
	IsRealtime bool
	Resp       *TransactionResp
}

type TransactionReq struct {
	StockQuery
	Offset uint16
	Size   uint16
}

type TransactionResp struct {
	Count uint16
	List  []TransactionRespItem
}

type TransactionRespItem struct {
	Tick      uint16 // hhmm
	Price     int64
	Volume    int64
	Number    int64
	BuyOrSell int64
	R0        int64
}

func (c *Client) Transaction(q StockQuery, isRealtime bool) (*TransactionResp, error) {
	t := &Transaction{
		Req: &TransactionReq{
			StockQuery: q,
		},
		IsRealtime: isRealtime,
	}
	t.Req.Size = 0x0384
	if isRealtime {
		t.Req.Size = 0x0032
	}
	t.SetDebug(c.ctx)

	err := do(c, c.dataConn, t)
	if err != nil {
		return nil, err
	}
	return t.Resp, nil
}

func (obj *Transaction) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.Method = 0x0FC5
	header.Type0 = 0x0100
	if obj.IsRealtime {
		header.Type0 = 0x1C00
	}
	return nil
}

func (obj *Transaction) MarshalReqBody(ctx context.Context) ([]byte, error) {
	var err error
	buf := bytes.NewBuffer(nil)
	err = binary.Write(buf, binary.LittleEndian, obj.Req.FixedSize())
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, obj.Req.Offset)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, obj.Req.Size)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (obj *Transaction) UnmarshalResp(ctx context.Context, data []byte) error {
	var err error
	obj.Resp = &TransactionResp{}
	cursor := 0
	obj.Resp.Count, err = ReadInt(data, &cursor, obj.Resp.Count)
	if err != nil {
		return err
	}
	priceOffset := int64(0)
	for i := 0; i < int(obj.Resp.Count); i += 1 {
		item := TransactionRespItem{}
		item.Tick, err = ReadInt(data, &cursor, item.Tick)
		if err != nil {
			return err
		}
		item.Price, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		item.Price += priceOffset
		priceOffset = item.Price

		item.Volume, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}

		item.Number, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}

		item.BuyOrSell, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		item.R0, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		obj.Resp.List = append(obj.Resp.List, item)
	}
	return nil
}
