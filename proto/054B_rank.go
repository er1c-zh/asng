package proto

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) Rank(schemaKey string) (*ListResp, error) {
	var err error
	r := Rank{}

	switch schemaKey {
	case "delta-desc-all":
		r.StaticCodec.SetContentHex(c.ctx, "06002e0000002a0001000500000001000000")
	case "delta-desc-exclude-bj":
		r.StaticCodec.SetContentHex(c.ctx, "06002e0000002a0000000500260001000000")
	case "region-desc-all":
		r.StaticCodec.SetContentHex(c.ctx, "0600000000002a0000000500040001000000")
	default:
		// all A share order by code asc
		r.StaticCodec.SetContentHex(c.ctx, "0600000000002a0000000500000001000000")
	}

	// r.SetDebug(c.ctx)

	err = do(c, c.dataConn, &r)
	if err != nil {
		return nil, err
	}
	return &r.Resp, nil
}

type Rank struct {
	BlankCodec
	StaticCodec
	Resp ListResp
}

func (r *Rank) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.Method = 0x054B
	header.PacketType = 2
	return nil
}

func (r *Rank) MarshalReqBody(ctx context.Context) ([]byte, error) {
	return r.StaticCodec.MarshalReqBody(ctx)
}

func (r *Rank) UnmarshalResp(ctx context.Context, data []byte) error {
	var err error
	cursor := 0
	r.Resp.Reserved0, err = ReadInt(data, &cursor, r.Resp.Reserved0)
	if err != nil {
		return err
	}
	r.Resp.Count, err = ReadInt(data, &cursor, r.Resp.Count)
	if err != nil {
		return err
	}
	defer func() {
		if r.Debug {
			j, _ := json.MarshalIndent(r.Resp, "", "  ")
			fmt.Println(string(j))
		}
	}()
	for i := 0; i < int(r.Resp.Count); i++ {
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
		item.Reserved0, err = ReadTDXInt(data, &cursor)
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

		r.Resp.List = append(r.Resp.List, item)
	}

	return nil
}
