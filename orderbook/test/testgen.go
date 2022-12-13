package test

import (
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/e-zhydzetski/strips-tt/orderbook"
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
func GenerateOrders(count int) []Order {
	rand.Seed(time.Now().UnixNano())

	var priceMean = 5000

	orders := make([]Order, count)
	for i := 0; i < count; i++ {
		order := Order{
			ID:    uuid.New().String(),
			Type:  orderbook.OTAsk,
			Value: orderbook.Value(1 + rand.Intn(int(maxValue))),
			Price: orderbook.PLMarket, // market price
		}
		if probability(50) {
			order.Type = orderbook.OTBid
		}
		if !probability(marketProb) {
			order.Price = orderbook.PriceLimit(1 + int(rand.NormFloat64()*float64(priceDev)+float64(priceMean)))
		}
		orders[i] = order

		if i%10 == 0 {
			priceMean = 1 + int(rand.NormFloat64()*float64(priceDev)+float64(priceMean))
		}
	}

	return orders
}

//nolint:gosec // unsecure random is ok
func probability(probInPercents uint8) bool {
	return uint8(rand.Intn(100)) < probInPercents
}
