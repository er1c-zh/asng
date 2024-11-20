package proto

import (
	"context"
)

type Heartbeat struct {
	BlankCodec
	StaticCodec
}

func (h *Heartbeat) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.Method = 0x000D
	header.PacketType = 1
	return h.StaticCodec.FillReqHeader(ctx, header)
}

func (h *Heartbeat) MarshalReqBody(ctx context.Context) ([]byte, error) {
	h.StaticCodec.SetContentHex(ctx, "01")
	return h.StaticCodec.MarshalReqBody(ctx)
}
