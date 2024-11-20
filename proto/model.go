package proto

import (
	"context"
	"encoding/hex"
	"fmt"
)

type ApiCode uint16

func (c ApiCode) String() string {
	return fmt.Sprintf("%04X", uint16(c))
}

type ReqHeader struct {
	MagicNumber uint8  // MagicNumber
	SeqID       uint32 // 请求编号
	PacketType  uint8
	PkgLen1     uint16
	PkgLen2     uint16
	Method      ApiCode // method 请求方法
}

type RespHeader struct {
	I1     uint32
	Flag   uint8
	SeqID  uint32 // 请求编号
	I3     uint8
	Method ApiCode // method
	// TODO 有时这个 PkgDataSize > RawDataSize
	PkgDataSize uint16 // 长度
	RawDataSize uint16 // 未压缩长度
}

type Codec interface {
	FillReqHeader(ctx context.Context, h *ReqHeader) error
	MarshalReqBody(ctx context.Context) ([]byte, error)
	UnmarshalResp(ctx context.Context, data []byte) error
	IsDebug(ctx context.Context) bool
	NeedEncrypt(ctx context.Context) bool
}

type BlankCodec struct {
	Debug   bool
	Encrypt bool
}

func (BlankCodec) MarshalReqBody(context.Context) ([]byte, error) {
	return nil, nil
}

func (BlankCodec) FillReqHeader(context.Context, *ReqHeader) error {
	return nil
}

func (BlankCodec) UnmarshalResp(context.Context, []byte) error {
	return nil
}

func (c *BlankCodec) IsDebug(ctx context.Context) bool {
	return c.Debug
}

func (c *BlankCodec) SetDebug(ctx context.Context) {
	c.Debug = true
}

func (c *BlankCodec) NeedEncrypt(ctx context.Context) bool {
	return false
}

func (c *BlankCodec) SetNeedEncrypt(context.Context, string) {
	c.Encrypt = true
}

type StaticCodec struct {
	BlankCodec
	ContentHex string
}

func (c *StaticCodec) SetContentHex(ctx context.Context, s string) {
	c.ContentHex = s
}

func (c *StaticCodec) MarshalReqBody(context.Context) ([]byte, error) {
	return hex.DecodeString(c.ContentHex)
}
