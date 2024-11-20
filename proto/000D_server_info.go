package proto

import (
	"context"
)

type ServerInfo struct {
	BlankCodec
	StaticCodec
	Resp ServerInfoResp
}

type ServerInfoResp struct {
	Name string
}

func (c *Client) ServerInfo() (*ServerInfo, error) {
	var err error
	s := ServerInfo{}
	s.SetContentHex(c.ctx, "01")
	err = do(c, c.dataConn, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *ServerInfo) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.Method = 0x000D
	header.PacketType = 1
	return nil
}

func (s *ServerInfo) MarshalReqBody(ctx context.Context) ([]byte, error) {
	return s.StaticCodec.MarshalReqBody(ctx)
}

func (s *ServerInfo) UnmarshalResp(ctx context.Context, data []byte) error {
	var err error
	cursor := 68
	s.Resp.Name, err = ReadTDXString(data, &cursor, 80)
	if err != nil {
		return err
	}
	return nil
}
