package proto

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"sort"
)

var RealtimeSubscribeType0 uint16 = 0x0029
var RealtimeSubscribeHandler RespHandler = func(c *Client, h *RespHeader, data []byte) {
	a := &RealtimeInfo{}
	err := a.UnmarshalResp(c.ctx, data)
	if err != nil {
		c.Log("RealtimeSubscribeHandler fail: %s", err.Error())
		return
	}
	for _, item := range a.Resp.ItemList {
		if c.RealtimeConsumer == nil {
			c.LogDebug("RealtimeSubscribeHandler handler miss")
			continue
		}
		c.RealtimeConsumer(item)
	}
}

func (c *Client) RegisterRealtimeConsumer(h RealtimeConsumer) {
	c.RealtimeConsumer = h
}

type RealtimeConsumer func(RealtimeInfoRespItem)

func (c *Client) RealtimeInfo(stock []StockQuery) (*RealtimeInfoResp, error) {
	return c.realtimeInfo(stock, false)
}

func (c *Client) RealtimeInfoSubscribe(stock []StockQuery) (*RealtimeInfoResp, error) {
	return c.realtimeInfo(stock, true)
}

func (c *Client) realtimeInfo(stock []StockQuery, subscribe bool) (*RealtimeInfoResp, error) {
	realtime := &RealtimeInfo{}

	realtime.SetDebug(c.ctx)
	realtime.Subscribe = subscribe

	realtime.Req = &RealtimeInfoReq{
		Size: uint16(len(stock)),
	}
	for _, s := range stock {
		realtime.Req.ItemList = append(realtime.Req.ItemList, RealtimeInfoReqItem{
			Market: s.Market,
			Code:   [10]byte{s.Code[0], s.Code[1], s.Code[2], s.Code[3], s.Code[4], s.Code[5]},
		})
	}

	err := do(c, c.dataConn, realtime)
	if err != nil {
		return nil, err
	}

	return realtime.Resp, nil
}

type RealtimeInfo struct {
	BlankCodec
	Subscribe bool
	Req       *RealtimeInfoReq
	Resp      *RealtimeInfoResp
}

type RealtimeInfoReq struct {
	Size     uint16
	ItemList []RealtimeInfoReqItem
}

type RealtimeInfoReqItem struct {
	Market uint8
	Code   [10]byte
}

type RealtimeInfoResp struct {
	Data     []byte
	Count    uint16
	ItemList []RealtimeInfoRespItem
}
type RealtimeInfoRespItem struct {
	StockQuery
	CurrentPrice        int64
	YesterdayCloseDelta int64
	OpenDelta           int64
	HighDelta           int64
	LowDelta            int64
	YesterdayClose      int64
	Open                int64
	High                int64
	Low                 int64
	TotalVolume         int64
	CurrentVolume       int64
	TotalAmount         float64

	TickNo           uint16
	TickInHHmmss     uint32 // time in in hhmmss
	AfterHoursVolume int64  // 盘后量
	SellAmount       int64
	BuyAmount        int64
	TickPriceDelta   int64
	OpenAmount       int64
	OrderBookRaw     [4 * 5]int64 // order book
	OrderBookRows    []OrderBookRow
	RUint0           uint16
	RUint1           uint32
	RUint2           uint32
	RB               string
	RIntArray2       [4*5 + 4]int64
}

type OrderBookRow struct {
	Price  int64
	Volume int64
}

func (obj *RealtimeInfoRespItem) Unmarshal(ctx context.Context, buf []byte, cursor *int) error {
	var err error
	obj.Market, err = ReadInt(buf, cursor, obj.Market)
	if err != nil {
		return err
	}
	obj.Code, err = ReadCode(buf, cursor)
	if err != nil {
		return err
	}
	obj.TickNo, err = ReadInt(buf, cursor, obj.TickNo)
	if err != nil {
		return err
	}
	obj.CurrentPrice, err = ReadTDXInt(buf, cursor)
	if err != nil {
		return err
	}
	obj.YesterdayCloseDelta, err = ReadTDXInt(buf, cursor)
	if err != nil {
		return err
	}
	obj.OpenDelta, err = ReadTDXInt(buf, cursor)
	if err != nil {
		return err

	}
	obj.HighDelta, err = ReadTDXInt(buf, cursor)
	if err != nil {
		return err
	}
	obj.LowDelta, err = ReadTDXInt(buf, cursor)
	if err != nil {
		return err
	}
	obj.TickInHHmmss, err = ReadInt(buf, cursor, obj.TickInHHmmss)
	if err != nil {
		return err
	}

	obj.AfterHoursVolume, err = ReadTDXInt(buf, cursor)
	if err != nil {
		return err
	}

	obj.TotalVolume, err = ReadTDXInt(buf, cursor)
	if err != nil {
		return err
	}
	obj.CurrentVolume, err = ReadTDXInt(buf, cursor)
	if err != nil {
		return err
	}

	obj.TotalAmount, err = ReadTDXFloat(buf, cursor)
	if err != nil {
		return err
	}

	obj.SellAmount, err = ReadTDXInt(buf, cursor)
	if err != nil {
		return err
	}
	obj.BuyAmount, err = ReadTDXInt(buf, cursor)
	if err != nil {
		return err
	}
	obj.TickPriceDelta, err = ReadTDXInt(buf, cursor)
	if err != nil {
		return err
	}
	obj.OpenAmount, err = ReadTDXInt(buf, cursor)
	if err != nil {
		return err
	}

	for i := 0; i < 4*5; i += 1 {
		obj.OrderBookRaw[i], err = ReadTDXInt(buf, cursor)
		if err != nil {
			return err
		}
	}

	tmp, err := ReadByteArray(buf, cursor, 10)
	if err != nil {
		return err
	}
	obj.RB = hex.Dump(tmp)
	// obj.RUint0, err = ReadInt(buf, cursor, obj.RUint0)
	// if err != nil {
	// 	return err
	// }
	// obj.RUint1, err = ReadInt(buf, cursor, obj.RUint1)
	// if err != nil {
	// 	return err
	// }
	// obj.RUint2, err = ReadInt(buf, cursor, obj.RUint2)
	// if err != nil {
	// 	return err
	// }

	for i := 0; i < 4*5+4; i += 1 {
		obj.RIntArray2[i], err = ReadTDXInt(buf, cursor)
		if err != nil {
			return err
		}
	}

	return nil
}

func (obj *RealtimeInfo) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.Method = 0x0547
	header.Type0 = 0x002A
	if obj.Subscribe {
		header.Type0 = 0x0029
	}
	return nil
}

func (obj *RealtimeInfo) MarshalReqBody(ctx context.Context) ([]byte, error) {
	var err error
	buf := bytes.NewBuffer(nil)
	err = binary.Write(buf, binary.LittleEndian, obj.Req.Size)
	if err != nil {
		return nil, err
	}
	for _, item := range obj.Req.ItemList {
		err = binary.Write(buf, binary.LittleEndian, item)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (obj *RealtimeInfo) UnmarshalResp(ctx context.Context, data []byte) error {
	data = decryptSimpleXOR(data, keySimpleXOR0547)
	obj.Resp = &RealtimeInfoResp{
		Data: data,
	}
	var err error
	cursor := 0
	obj.Resp.Count, err = ReadInt(data, &cursor, obj.Resp.Count)
	if err != nil {
		return err
	}

	for i := 0; i < int(obj.Resp.Count); i += 1 {
		item := RealtimeInfoRespItem{}
		err = item.Unmarshal(ctx, data, &cursor)
		if err != nil {
			return err
		}

		item.Open = item.CurrentPrice + item.OpenDelta
		item.YesterdayClose = item.CurrentPrice + item.YesterdayCloseDelta
		item.High = item.CurrentPrice + item.HighDelta
		item.Low = item.CurrentPrice + item.LowDelta

		for i := 0; i < len(item.OrderBookRaw); i += 4 {
			item.OrderBookRows = append(item.OrderBookRows, OrderBookRow{
				Price:  item.CurrentPrice + item.OrderBookRaw[i],
				Volume: item.OrderBookRaw[i+2],
			})
			item.OrderBookRows = append(item.OrderBookRows, OrderBookRow{
				Price:  item.CurrentPrice + item.OrderBookRaw[i+1],
				Volume: item.OrderBookRaw[i+3],
			})
		}
		sort.Slice(item.OrderBookRows, func(i, j int) bool {
			return item.OrderBookRows[i].Price < item.OrderBookRows[j].Price
		})

		obj.Resp.ItemList = append(obj.Resp.ItemList, item)
	}

	return nil
}
