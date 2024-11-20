package proto

import (
	"bytes"
	"context"
	"encoding/binary"
)

type MetaDescMap struct {
	BlankCodec
	Req  MetaDescMapReq
	Resp *MetaDescMapResp
}

type MetaDescMapReq struct {
	Offset   uint32
	PageSize uint16
}

type MetaDescMapResp struct {
	Reserved0 uint32
	Count     uint16
	List      []MetaDesc
}

type MetaDesc struct {
	Reserved0 []byte // 5
	ID        string // 9 ascii
	Desc      string // 28 gbk
	Reserved1 []byte // 22
}

func (c *Client) MetaDescMap(conn *ConnRuntime, offset uint32) (*MetaDescMapResp, error) {
	var err error
	m := &MetaDescMap{}

	m.Req.Offset = offset
	m.Req.PageSize = 500

	err = do(c, conn, m)
	if err != nil {
		return nil, err
	}
	return m.Resp, nil
}

func (m *MetaDescMap) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.MagicNumber = 1
	header.Method = 0x23F5
	header.PacketType = 1
	return nil
}

func (m *MetaDescMap) MarshalReqBody(ctx context.Context) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := binary.Write(buf, binary.LittleEndian, m.Req)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *MetaDescMap) UnmarshalResp(ctx context.Context, data []byte) error {
	var err error
	m.Resp = &MetaDescMapResp{}
	resp := m.Resp
	cursor := 0
	resp.Reserved0, err = ReadInt(data, &cursor, resp.Reserved0)
	if err != nil {
		return err
	}
	resp.Count, err = ReadInt(data, &cursor, resp.Count)
	if err != nil {
		return err
	}

	for i := 0; i < int(resp.Count); i += 1 {
		item := MetaDesc{}
		item.Reserved0, err = ReadByteArray(data, &cursor, 5)
		if err != nil {
			return err
		}
		item.ID, err = ReadTDXString(data, &cursor, 9)
		if err != nil {
			return err
		}
		item.Desc, err = ReadTDXString(data, &cursor, 28)
		if err != nil {
			return err
		}
		item.Reserved1, err = ReadByteArray(data, &cursor, 22)
		if err != nil {
			return err
		}

		resp.List = append(resp.List, item)
	}

	return nil
}
