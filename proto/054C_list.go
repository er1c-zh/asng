package proto

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

func (c *Client) List(stockList []StockQuery) (*ListResp, error) {
	var err error
	l := List{}
	l.SetContentHex(c.ctx, "0500000000000000"+fmt.Sprintf("%02X", len(stockList))+"00")
	l.ListReq.Items = make([]ListReqItem, 0, len(stockList))
	for _, stock := range stockList {
		reqItem := ListReqItem{
			Market: stock.Market,
		}
		reqItem.Code, err = GenerateCodeBytesArray(stock.Code)
		if err != nil {
			return nil, err
		}

		l.ListReq.Items = append(l.ListReq.Items, reqItem)
	}
	err = do(c, c.dataConn, &l)
	if err != nil {
		return nil, err
	}
	return &l.Resp, nil
}

type List struct {
	BlankCodec
	ListReq
	Resp ListResp
}

type ListReq struct {
	StaticCodec
	Items []ListReqItem
}

type ListReqItem struct {
	Market uint8
	Code   [6]byte
}

type ListResp struct {
	Reserved0 uint16
	Count     uint16
	List      []ListRespItem
}
type ListRespItem struct {
	Market             uint8
	Code               string
	SplitFlag          uint16
	CurPrice           int64
	YesterdayOpenDelta int64
	OpenDelta          int64
	HighDelta          int64
	LowDelta           int64
	Reserved0          int64
	NegativeCurPrice   int64 // ?
	TotalVolume        int64
	CurrentVolume      int64
	TotalAmount        float64
	SellVolume         int64
	BuyVolume          int64
	Reserved1          int64
	Reserved2          int64

	BuyPriceDelta1  int64
	SellPriceDelta1 int64
	BuyVolume1      int64
	SellVolume1     int64

	Reserved3 []byte

	SplitFlagEnd uint16
}

func (l *List) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.Method = 0x054C
	header.PacketType = 2
	return nil
}

func (l *List) UnmarshalResp(ctx context.Context, data []byte) error {
	var err error
	cursor := 0
	l.Resp.Reserved0, err = ReadInt(data, &cursor, l.Resp.Reserved0)
	if err != nil {
		return err
	}
	l.Resp.Count, err = ReadInt(data, &cursor, l.Resp.Count)
	if err != nil {
		return err
	}
	defer func() {
		j, _ := json.MarshalIndent(l.Resp, "", "  ")
		fmt.Println(string(j))
	}()
	for i := 0; i < int(l.Resp.Count); i++ {
		item := ListRespItem{}
		item.Market, err = ReadInt(data, &cursor, item.Market)
		if err != nil {
			return err
		}
		item.Code, err = ReadCode(data, &cursor)
		if err != nil {
			return err
		}
		item.SplitFlag, err = ReadInt(data, &cursor, item.SplitFlag)
		if err != nil {
			return err
		}
		item.CurPrice, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		item.YesterdayOpenDelta, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		item.OpenDelta, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		item.HighDelta, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		item.LowDelta, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		_, err = ReadInt(data, &cursor, uint32(0))
		if err != nil {
			return err
		}
		item.NegativeCurPrice, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		item.TotalVolume, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		item.CurrentVolume, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		item.TotalAmount, err = ReadTDXFloat(data, &cursor)
		if err != nil {
			return err
		}
		item.SellVolume, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		item.BuyVolume, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		item.Reserved1, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		item.Reserved2, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}

		item.BuyPriceDelta1, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		item.SellPriceDelta1, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		item.BuyVolume1, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}
		item.SellVolume1, err = ReadTDXInt(data, &cursor)
		if err != nil {
			return err
		}

		item.Reserved3, err = ReadByteArray(data, &cursor, 10+12+8+24)
		if err != nil {
			return err
		}
		item.SplitFlagEnd, err = ReadInt(data, &cursor, item.SplitFlagEnd)
		if err != nil {
			return err
		}

		l.Resp.List = append(l.Resp.List, item)
	}

	return nil
}

func (l *List) MarshalReqBody(ctx context.Context) ([]byte, error) {
	optionData, err := l.StaticCodec.MarshalReqBody(ctx)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(optionData)
	err = binary.Write(buf, binary.LittleEndian, l.Items)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
