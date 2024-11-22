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

type HeaderIdentifer struct {
	SeqID uint8 // 请求编号
	Group uint8
	Type0 uint16 // 方法分类？
}

func (h HeaderIdentifer) UniqKey() string {
	return fmt.Sprintf("%02X%02X%04X", h.SeqID, h.Group, h.Type0)
}

type ReqHeader struct {
	MagicNumber uint8 // MagicNumber
	HeaderIdentifer
	PacketType uint8
	PkgLen1    uint16
	PkgLen2    uint16
	Method     ApiCode // method 请求方法
}

type RespHeader struct {
	R0   uint32
	Flag uint8
	HeaderIdentifer
	R1     uint8
	Method ApiCode // method
	// 有时这个 PkgDataSize > RawDataSize 不根据这个判断是否是压缩处理过
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
	RespHex string
}

func (BlankCodec) MarshalReqBody(context.Context) ([]byte, error) {
	return nil, nil
}

func (BlankCodec) FillReqHeader(context.Context, *ReqHeader) error {
	return nil
}

func (c *BlankCodec) UnmarshalResp(ctx context.Context, data []byte) error {
	c.RespHex = hex.Dump(data)
	return nil
}

func (c *BlankCodec) IsDebug(ctx context.Context) bool {
	return c.Debug
}

func (c *BlankCodec) SetDebug(ctx context.Context) {
	c.Debug = true
}

func (c *BlankCodec) NeedEncrypt(ctx context.Context) bool {
	return c.Encrypt
}

func (c *BlankCodec) SetNeedEncrypt(context.Context) {
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
