package proto

import (
	"bytes"
	"context"
	"encoding/binary"
)

func (c *Client) Realtime(stock []StockQuery) (*RealtimeResp, error) {
	realtime := &Realtime{}

	realtime.SetDebug(c.ctx)

	realtime.Req = &RealtimeReq{
		Size: uint16(len(stock)),
	}
	for _, s := range stock {
		realtime.Req.ItemList = append(realtime.Req.ItemList, RealtimeReqItem{
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

type Realtime struct {
	BlankCodec
	Req  *RealtimeReq
	Resp *RealtimeResp
}

type RealtimeReq struct {
	Size     uint16
	ItemList []RealtimeReqItem
}

type RealtimeReqItem struct {
	Market uint8
	Code   [10]byte
}

type RealtimeResp struct {
	Data     []byte
	Count    uint16
	ItemList []RealtimeRespItem
}
type RealtimeRespItem struct {
	Market              uint8
	Code                string
	CurrentPrice        int64
	YesterdayCloseDelta int64
	OpenDelta           int64
	HighDelta           int64
	LowDelta            int64
	TotalVolume         int64
	CurrentVolume       int64
	TotalAmount         float64
}

func (obj *RealtimeRespItem) Unmarshal(ctx context.Context, buf []byte, cursor *int) error {
	var err error
	obj.Market, err = ReadInt(buf, cursor, obj.Market)
	if err != nil {
		return err
	}
	obj.Code, err = ReadCode(buf, cursor)
	if err != nil {
		return err
	}
	_, err = ReadByteArray(buf, cursor, 2)
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
	_, err = ReadInt(buf, cursor, uint32(0))
	if err != nil {
		return err
	}

	_, err = ReadTDXInt(buf, cursor)
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

	for i := 0; i < 4; i += 1 {
		_, err = ReadTDXInt(buf, cursor)
		if err != nil {
			return err
		}
	}

	for i := 0; i < 4*5; i += 1 {
		_, err = ReadTDXInt(buf, cursor)
		if err != nil {
			return err
		}
	}

	_, err = ReadByteArray(buf, cursor, 2+4+4)
	if err != nil {
		return err
	}

	for i := 0; i < 4*5+4; i += 1 {
		_, err = ReadTDXInt(buf, cursor)
		if err != nil {
			return err
		}
	}

	return nil
}

func (obj *Realtime) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.Method = 0x0547
	return nil
}

func (obj *Realtime) MarshalReqBody(ctx context.Context) ([]byte, error) {
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

func (obj *Realtime) UnmarshalResp(ctx context.Context, data []byte) error {
	data = decryptSimpleXOR(data, keySimpleXOR0547)
	obj.Resp = &RealtimeResp{
		Data: data,
	}
	var err error
	cursor := 0
	obj.Resp.Count, err = ReadInt(data, &cursor, obj.Resp.Count)
	if err != nil {
		return err
	}

	for i := 0; i < int(obj.Resp.Count); i += 1 {
		item := RealtimeRespItem{}
		err = item.Unmarshal(ctx, data, &cursor)
		if err != nil {
			return err
		}
		obj.Resp.ItemList = append(obj.Resp.ItemList, item)
	}

	return nil
}
