package test

import (
	"testing"

	"github.com/e-zhydzetski/strips-tt/orderbook"
)

func BenchmarkOrderbook(b *testing.B) {
	orderBook := orderbook.NewOrderBook()
	orders := GenerateOrders(1000)

	idx := 0
	nextOrder := func() Order {
		idx = (idx + 1) % len(orders)
		return orders[idx]
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		order := nextOrder()
		b.StartTimer()
		switch order.Type {
		case orderbook.OTBid:
			orderBook.Bid(order.ID, order.Value, order.Price)
		case orderbook.OTAsk:
			orderBook.Ask(order.ID, order.Value, order.Price)
		}
	}
}
