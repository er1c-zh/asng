package proto

import (
	"asng/models"
	"bytes"
	"context"
	"encoding/binary"
	"slices"
	"time"
)

type Transaction struct {
	BlankCodec
	Req        *TransactionReq
	IsRealtime bool
	IsHistory  bool
	Resp       *TransactionResp
}

type TransactionReq struct {
	models.StockIdentity
	Date   time.Time
	Offset uint16
	Size   uint16
}

type TransactionResp struct {
	R0    uint32
	Count uint16
	List  []TransactionRespItem
}

type TransactionRespItem struct {
	Tick      uint16 // minutes from today 00:00
	Price     int64
	Volume    int64
	Number    int64
	BuyOrSell int64
	R0        int64
}

func (c *Client) TXToday(q models.StockIdentity) ([]TransactionRespItem, error) {
	offset := uint16(0)
	result := make([][]TransactionRespItem, 0, 10)
	for {
		resp, err := c.TXRealtime(q, offset)
		if err != nil {
			return nil, err
		}

		result = append(result, resp.List)
		offset += resp.Count

		if resp.Count < TXHistoryBatchSize {
			break
		}
	}
	slices.Reverse(result)
	return slices.Concat(result...), nil
}

func (c *Client) TXRealtime(q models.StockIdentity, offset uint16) (*TransactionResp, error) {
	t := &Transaction{
		Req: &TransactionReq{
			StockIdentity: q,
			// Size:          0x0032, // 50
			Size:   TXHistoryBatchSize,
			Offset: offset,
		},
		IsRealtime: true,
	}
	err := do(c, c.dataConn, t)
	if err != nil {
		return nil, err
	}
	return t.Resp, nil
}

const (
	TXHistoryBatchSize = 0x384
)

func (c *Client) TXHistoryOneDay(q models.StockIdentity, date time.Time) ([]TransactionRespItem, error) {
	offset := uint16(0)
	result := make([][]TransactionRespItem, 0, 10)
	for {
		resp, err := c.TXHistory(q, date, offset)
		if err != nil {
			return nil, err
		}

		result = append(result, resp.List)
		offset += resp.Count

		if resp.Count < TXHistoryBatchSize {
			break
		}
	}
	slices.Reverse(result)
	return slices.Concat(result...), nil
}

func (c *Client) TXHistory(q models.StockIdentity, date time.Time, offset uint16) (*TransactionResp, error) {
	t := &Transaction{
		Req: &TransactionReq{
			StockIdentity: q,
			Date:          date,
			Offset:        offset,
			Size:          TXHistoryBatchSize,
		},
		IsHistory: true,
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
	if obj.IsHistory {
		header.Method = 0x0FC6
		// header.Method = 0x0FC3 // L2
	}
	if obj.IsRealtime {
		// header.Type0 = 0x0100
		header.Type0 = 0x1C00
	} else if obj.IsHistory {
		header.Type0 = 0x0101
	}
	return nil
}

func (obj *Transaction) MarshalReqBody(ctx context.Context) ([]byte, error) {
	var err error
	buf := bytes.NewBuffer(nil)
	if obj.IsHistory {
		err = binary.Write(buf, binary.LittleEndian, FormatTDXDate(obj.Req.Date, TDXTimeTypeYYYYMMDD))
		if err != nil {
			return nil, err
		}
	}
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
	if obj.IsHistory {
		obj.Resp.R0, err = ReadInt(data, &cursor, obj.Resp.R0)
		if err != nil {
			return err
		}
	}
	obj.Resp.List = make([]TransactionRespItem, 0, obj.Resp.Count)
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
