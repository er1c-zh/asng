package proto

import (
	"asng/models"
	"bytes"
	"context"
	"encoding/binary"
	"regexp"
	"strings"
	"sync"

	"github.com/mozillazg/go-slugify"
)

var (
	slugOnce  sync.Once
	slugRegex = regexp.MustCompile(`[a-z]+`)
)

func getPinYinInitial(s string) string {
	slugOnce.Do(func() {})
	b := bytes.NewBuffer(nil)
	for _, s := range strings.Split(slugify.Slugify(s), slugify.Separator) {
		if slugRegex.MatchString(s) {
			b.WriteByte(s[0])
		} else {
			b.WriteString(s)
		}
	}
	return b.String()
}

func (c *Client) StockMetaAll() (*models.StockMetaAll, error) {
	d := models.StockMetaAll{}
	for _, market := range []models.MarketType{models.MarketSZ, models.MarketSH, models.MarketBJ} {
		cursor := uint32(0)
		for {
			resp, err := c.StockMeta(market, cursor)
			if err != nil {
				return nil, err
			}
			for _, item := range resp.List {
				d.StockList = append(d.StockList, models.StockMetaItem{
					Code:           item.Code,
					Market:         market,
					Desc:           item.Desc,
					PinYinInitial:  getPinYinInitial(item.Desc),
					Scale:          item.Scale,
					F0:             item.F0,
					YesterdayClose: item.YesterdayClose,
					BlockID:        item.BlockID,
					F3:             item.F3,
				})
			}
			cursor += uint32(len(resp.List))
			if resp.Count == 0 {
				break
			}
		}
	}
	return &d, nil
}

func (c *Client) StockMeta(market models.MarketType, offset uint32) (*StockMetaResp, error) {
	var err error
	s := StockMeta{}
	s.Req = &StockMetaReq{
		Market: market,
		Size:   0x0640,
		Offset: offset,
		R1:     0,
	}
	err = do(c, c.dataConn, &s)
	if err != nil {
		return nil, err
	}
	return s.Resp, nil
}

type StockMeta struct {
	BlankCodec
	Req  *StockMetaReq
	Resp *StockMetaResp
}

/*
0000000000 0c 16 18 6e 00 01 10 00 10 00 4d 04 00 00 00 00   ...n......M.....
0000000010 00 00 40 06 00 00 00 00 00 00                     ..@.......
*/
type StockMetaReq struct {
	Market models.MarketType
	Offset uint32
	Size   uint32
	R1     uint32
}

type StockMetaResp struct {
	Count uint16
	List  []StockMetaItem
}

type StockMetaItem struct {
	Code           string
	Scale          uint16
	Desc           string
	F0             float64
	YesterdayClose float64
	BlockID        uint16
	F3             uint16
}

func (StockMeta) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.MagicNumber = 0x0C
	header.Type0 = 0x006E
	header.PacketType = 1
	header.Method = 0x044D
	return nil
}

func (obj *StockMeta) MarshalReqBody(ctx context.Context) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := binary.Write(buf, binary.LittleEndian, obj.Req)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

/*
unsigned __int16 __cdecl parseSingle044Dsub_81CBE0(char *Src, char *a2)

	{
	  bool v2; // zf
	  unsigned __int16 result; // ax
	  const char *TdxPYStr; // eax

	  memset(a2, 0, 0x168u);
	  memmove(a2, Src, 6u); // code
	  a2[6] = 0;
	  memmove(a2 + 31, Src + 8, 0x10u); // gbk desc
	  a2[47] = 0;
	  v2 = a2[329] == 0;
	  *(float *)(a2 + 0x4E) = (float)*((__int16 *)Src + 3); // scale
	  a2[76] = Src[28];
	  *(float *)(a2 + 90) = *((float *)Src + 6);
	  *((float *)a2 + 69) = *(float *)(Src + 29);
	  result = *(_WORD *)(Src + 33);
	  *((_WORD *)a2 + 136) = result;
	  *((_WORD *)a2 + 137) = *(_WORD *)(Src + 35);
	  if ( v2 )
	  {
	    TdxPYStr = (const char *)GetTdxPYStr(a2 + 31, 8);
	    return (unsigned __int16)strncpy(a2 + 329, TdxPYStr, 8u);
	  }
	  return result;
	}
*/
func (obj *StockMeta) UnmarshalResp(ctx context.Context, data []byte) error {
	var err error
	cursor := 0
	obj.Resp = &StockMetaResp{}
	obj.Resp.Count, err = ReadInt(data, &cursor, obj.Resp.Count)
	if err != nil {
		return err
	}

	for i := 0; i < int(obj.Resp.Count); i++ {
		item := StockMetaItem{}

		item.Code, err = ReadCode(data, &cursor)
		if err != nil {
			return err
		}
		item.Scale, err = ReadInt(data, &cursor, item.Scale)
		if err != nil {
			return err
		}
		item.Desc, err = ReadTDXString(data, &cursor, 16)
		if err != nil {
			return err
		}

		// _, err = ReadByteArray(data, &cursor, c0+37 /* item data length */ -cursor)
		// if err != nil {
		// 	return err
		// }

		item.F0, err = ReadTDXFloat(data, &cursor)
		if err != nil {
			return err
		}

		cursor += 1 // FIXME: really padding

		item.YesterdayClose, err = ReadTDXFloat(data, &cursor)
		if err != nil {
			return err
		}

		item.BlockID, err = ReadInt(data, &cursor, item.BlockID)
		if err != nil {
			return err
		}

		item.F3, err = ReadInt(data, &cursor, item.F3)
		if err != nil {
			return err
		}

		obj.Resp.List = append(obj.Resp.List, item)
	}
	return nil
}
