package proto

import (
	"context"
)

type MetaHandshake struct {
	BlankCodec
	StaticCodec
}

func (c *Client) MetaShakehand(conn *ConnRuntime) error {
	handshake := &MetaHandshake{}
	handshake.SetContentHex(c.ctx, "e5bb1c2fafe525941f32c6e5d53dfb415b734cc9cdbf0ac92021bfdd1eb06d2266ac09ae75259234e4b734ac2663ad911f32c6e5d53dfb411f32c6e5d53dfb4124288302ea5fc63df443512c51847720")
	err := do(c, conn, handshake)
	if err != nil {
		return err
	}
	return nil
}

func (mh *MetaHandshake) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.MagicNumber = 01
	header.Method = 0x2454
	header.PacketType = 1
	return nil
}

func (mh *MetaHandshake) MarshalReqBody(ctx context.Context) ([]byte, error) {
	return mh.StaticCodec.MarshalReqBody(ctx)
}
