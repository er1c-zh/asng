package proto

import (
	"context"
)

func (c *Client) Heartbeat() error {
	h := &Heartbeat{}
	h.SetDebug(c.ctx)
	err := do(c, c.dataConn, h)
	if err != nil {
		return err
	}
	return nil
}

type Heartbeat struct {
	BlankCodec
}

func (h *Heartbeat) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.Type0 = 0x0028
	header.Method = 0x0004
	header.PacketType = 0x02
	return nil
}
