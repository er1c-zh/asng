package models

import "fmt"

type StockIdentity struct {
	MarketType MarketType
	Code       string
}
type StockIdentityFixedSize struct {
	MarketType MarketType
	Code       [6]byte
}

type StockIdentityFixedSizeWithUint8MarketType struct {
	MarketType uint8
	Code       [6]byte
}

func (s StockIdentity) String() string {
	return fmt.Sprintf("%d%s", s.MarketType, s.Code)
}

func (s StockIdentity) CodeArray() [6]byte {
	if len(s.Code) != 6 {
		return [6]byte{}
	}
	return [6]byte{s.Code[0], s.Code[1], s.Code[2], s.Code[3], s.Code[4], s.Code[5]}
}

func (s StockIdentity) FixedSize() StockIdentityFixedSize {
	return StockIdentityFixedSize{
		MarketType: s.MarketType,
		Code:       s.CodeArray(),
	}
}

func (s StockIdentity) FixedSizeWithUint8MarketType() StockIdentityFixedSizeWithUint8MarketType {
	return StockIdentityFixedSizeWithUint8MarketType{
		MarketType: uint8(s.MarketType),
		Code:       s.CodeArray(),
	}
}
