package proto

import (
	"fmt"
)

type StockQuery struct {
	Market uint8
	Code   string
}

func (sq StockQuery) String() string {
	return fmt.Sprintf("%d%s", sq.Market, sq.Code)
}
