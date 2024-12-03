package proto

import (
	"asng/models"
	"fmt"
)

type StockQuery struct {
	Market uint8
	Code   string
}
type StockQueryFixedSize struct {
	Market uint16
	Code   [6]byte
}

func (sq StockQuery) String() string {
	return fmt.Sprintf("%d%s", sq.Market, sq.Code)
}

func (sq StockQuery) FixedSize() StockQueryFixedSize {
	return StockQueryFixedSize{
		Market: uint16(sq.Market),
		Code:   [6]byte{sq.Code[0], sq.Code[1], sq.Code[2], sq.Code[3], sq.Code[4], sq.Code[5]},
	}
}

func (sq StockQuery) ToID() models.StockIdentity {
	return models.StockIdentity{
		MarketType: models.MarketType(sq.Market),
		Code:       sq.Code,
	}
}
