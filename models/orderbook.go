package models

import (
	"slices"
)

type Orderbook struct {
	Asks []Order
	Bids []Order
}

func (o *Orderbook) Format() {
	slices.SortFunc(o.Asks, func(a, b Order) int {
		return int(a.Price - b.Price) // asc
	})
	slices.SortFunc(o.Bids, func(a, b Order) int {
		return int(b.Price - a.Price) // desc
	})
}

type Order struct {
	Price  int64
	Volume int64
}
