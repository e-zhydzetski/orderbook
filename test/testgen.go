package test

import (
	"math/rand"
	"strconv"

	"github.com/e-zhydzetski/orderbook"
)

type Order struct {
	ID    string
	Type  orderbook.OrderType
	Value orderbook.Value
	Price orderbook.PriceLimit
}

const (
	marketProb = 30
	priceDev   = 10
	maxValue   = orderbook.Value(100)
)

//nolint:gosec // unsecure random is ok
func NewOrdersGenerator(seed int64) func() Order {
	r := rand.New(rand.NewSource(seed))

	probability := func(probInPercents uint8) bool {
		return uint8(r.Intn(100)) < probInPercents
	}

	priceMean := 5000
	i := 0

	return func() Order {
		i++

		order := Order{
			ID:    strconv.Itoa(i),
			Type:  orderbook.OTAsk, // default
			Value: orderbook.Value(1 + rand.Intn(int(maxValue))),
			Price: orderbook.PLMarket, // market price by default
		}
		if probability(50) {
			order.Type = orderbook.OTBid
		}
		if !probability(marketProb) {
			order.Price = orderbook.PriceLimit(1 + int(rand.NormFloat64()*float64(priceDev)+float64(priceMean)))
		}

		if i%10 == 0 {
			priceMean = 1 + int(rand.NormFloat64()*float64(priceDev)+float64(priceMean))
		}

		return order
	}
}
