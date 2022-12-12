package orderbook

import (
	"time"
)

type Value uint64

type Price uint64

type OrderType byte

const (
	OTBid OrderType = iota
	OTAsk
)

type LimitOrder struct {
	ID         string
	Type       OrderType
	Value      Value
	Price      Price
	AcceptTime time.Time
}

type MarketOrder struct {
	ID         string
	Type       OrderType
	Value      Value
	AcceptTime time.Time
}
