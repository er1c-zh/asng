package proto

import "context"

func (c *Client) Nop(method ApiCode, type0 uint16, hex string) error {
	n := &Nop{
		method: method,
		type0:  type0,
		hex:    hex,
	}
	n.SetDebug(c.ctx)
	err := do(c, c.dataConn, n)
	if err != nil {
		return err
	}
	return nil
}

type Nop struct {
	BlankCodec
	StaticCodec
	method ApiCode
	type0  uint16
	hex    string
}

func (h *Nop) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.Type0 = h.type0
	header.Method = h.method
	header.PacketType = 1
	return nil
}

func (h *Nop) MarshalReqBody(ctx context.Context) ([]byte, error) {
	h.SetContentHex(ctx, h.hex)
	return h.StaticCodec.MarshalReqBody(ctx)
}
