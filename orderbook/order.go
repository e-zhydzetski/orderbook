package orderbook

import (
	"time"

	"github.com/e-zhydzetski/strips-tt/orderbook/memtable"
	"github.com/e-zhydzetski/strips-tt/orderbook/queue"
)

type Value uint64

type PriceLimit uint64

func (p PriceLimit) IsMarket() bool {
	return p == PLMarket
}

const PLMarket = PriceLimit(0)

type OrderType byte

const (
	OTBid OrderType = iota
	OTAsk
)

type Order struct {
	ID         string
	Type       OrderType
	Value      Value
	Price      PriceLimit
	AcceptTime time.Time
}

// NewOrderGroupFactoryFunc returns order group factory with common queue factory
func NewOrderGroupFactoryFunc() func() OrderGroup {
	queueFactory := queue.NewFactoryFunc[Order]()
	return func() OrderGroup {
		return OrderGroup{
			TotalValue: 0,
			Orders:     queueFactory(),
		}
	}
}

type OrderGroup struct {
	TotalValue Value
	Orders     memtable.Queue[Order]
}
