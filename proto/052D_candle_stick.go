package proto

import (
	"asng/models"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
)

func (c *Client) CandleStick(market models.MarketType, code string,
	periodType CandleStickPeriodType, cursor uint16) (*CandleStickResp, error) {
	var err error
	candleStick := &CandleStick{}
	candleStick.SetDebug(c.ctx)
	candleStick.Req = &CandleStickReq{
		MarketType: market,
		Code:       [6]byte{code[0], code[1], code[2], code[3], code[4], code[5]},
		Type:       periodType,
		Unit:       0x0001,
		Offset:     cursor,
		Size:       0x01A4,
	}
	err = do(c, c.dataConn, candleStick)
	if err != nil {
		return nil, err
	}
	candleStick.Resp.Cursor = cursor + candleStick.Resp.Size
	return candleStick.Resp, nil
}

type CandleStick struct {
	BlankCodec
	Req  *CandleStickReq
	Resp *CandleStickResp
}

type CandleStickPeriodType uint16

const (
	CandleStickPeriodType_5Min  CandleStickPeriodType = 0
	CandleStickPeriodType_15Min CandleStickPeriodType = 1
	CandleStickPeriodType_30Min CandleStickPeriodType = 2
	CandleStickPeriodType_1Hour CandleStickPeriodType = 3
	CandleStickPeriodType_Day   CandleStickPeriodType = 4
	CandleStickPeriodType_Week  CandleStickPeriodType = 5
	CandleStickPeriodType_Month CandleStickPeriodType = 6
	CandleStickPeriodType_1Min  CandleStickPeriodType = 7
)

type CandleStickReq struct {
	MarketType models.MarketType
	Code       [6]byte
	Type       CandleStickPeriodType
	Unit       uint16
	Offset     uint16
	Size       uint16
	Padding    [10]byte
}

type CandleStickResp struct {
	Size     uint16
	ItemList []CandleStickItem
	Cursor   uint16
}

type CandleStickItem struct {
	Open   int64
	Close  int64
	High   int64
	Low    int64
	Vol    float64 // 股数
	Amount float64 // 元

	TimeDesc string
}

func (c *CandleStick) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.Method = 0x052D
	return nil
}

func (c *CandleStick) MarshalReqBody(ctx context.Context) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := binary.Write(buf, binary.LittleEndian, c.Req)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *CandleStick) UnmarshalResp(ctx context.Context, body []byte) error {
	var err error
	c.Resp = &CandleStickResp{}

	cursor := 0
	c.Resp.Size, err = ReadInt(body, &cursor, c.Resp.Size)
	if err != nil {
		return err
	}
	priceDelta := int64(0)
	for i := 0; i < int(c.Resp.Size); i += 1 {
		item := CandleStickItem{}
		y, m, d, h, M, err := ReadTDXTime(body, &cursor, c.Req.Type)
		if err != nil {
			return err
		}
		item.TimeDesc = fmt.Sprintf("%04d-%02d-%02d %02d:%02d", y, m, d, h, M)

		nextPriceDelta := priceDelta
		item.Open, err = ReadTDXInt(body, &cursor)
		if err != nil {
			return err
		}
		nextPriceDelta += item.Open
		item.Open += priceDelta

		item.Close, err = ReadTDXInt(body, &cursor)
		if err != nil {
			return err
		}
		nextPriceDelta += (item.Close)
		item.Close += item.Open

		item.High, err = ReadTDXInt(body, &cursor)
		if err != nil {
			return err
		}
		item.High += item.Open

		item.Low, err = ReadTDXInt(body, &cursor)
		if err != nil {
			return err
		}
		item.Low += item.Open

		item.Vol, err = ReadTDXFloat(body, &cursor)
		if err != nil {
			return err
		}
		item.Amount, err = ReadTDXFloat(body, &cursor)
		if err != nil {
			return err
		}
		// item.UpCount, err = ReadInt(body, &cursor, item.UpCount)
		// if err != nil {
		// 	return err
		// }
		// item.DownCount, err = ReadInt(body, &cursor, item.DownCount)
		// if err != nil {
		// 	return err
		// }

		priceDelta = nextPriceDelta

		c.Resp.ItemList = append(c.Resp.ItemList, item)
	}

	return nil
}
