package proto

import (
	"context"
)

func (c *Client) ServerInfo() (*ServerInfoResp, error) {
	var err error
	s := ServerInfo{}
	s.SetContentHex(c.ctx, "01")
	err = do(c, c.dataConn, &s)
	if err != nil {
		return nil, err
	}
	return s.Resp, nil
}

type ServerInfo struct {
	BlankCodec
	StaticCodec
	Resp *ServerInfoResp
}

type ServerInfoResp struct {
	Name string
}

func (s *ServerInfo) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.Type0 = 0x0094
	header.Method = 0x000D
	header.PacketType = 1
	return nil
}

func (s *ServerInfo) MarshalReqBody(ctx context.Context) ([]byte, error) {
	return s.StaticCodec.MarshalReqBody(ctx)
}

func (s *ServerInfo) UnmarshalResp(ctx context.Context, data []byte) error {
	var err error
	s.Resp = &ServerInfoResp{}
	cursor := 68
	s.Resp.Name, err = ReadTDXString(data, &cursor, 80)
	if err != nil {
		return err
	}
	return nil
}
