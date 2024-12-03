package models

import "fmt"

type StockIdentity struct {
	MarketType MarketType
	Code       string
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
