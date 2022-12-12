package orderbook

import (
	"github.com/e-zhydzetski/strips-tt/orderbook/queue"
	"github.com/e-zhydzetski/strips-tt/orderbook/tree"
	"strings"
	"time"
)

// The result will be 0 if a == b, -1 if a < b, and +1 if a > b.
func lowToHighPrice(a LimitOrder, b LimitOrder) int {
	if a.Price < b.Price {
		return -1
	}
	if a.Price > b.Price {
		return 1
	}
	if a.AcceptTime.Before(b.AcceptTime) {
		return -1
	}
	if a.AcceptTime.After(b.AcceptTime) {
		return 1
	}
	return strings.Compare(a.ID, b.ID)
}

// The result will be 0 if a == b, -1 if a < b, and +1 if a > b.
func highToLowPrice(a LimitOrder, b LimitOrder) int {
	if a.Price > b.Price {
		return -1
	}
	if a.Price < b.Price {
		return 1
	}
	if a.AcceptTime.Before(b.AcceptTime) {
		return -1
	}
	if a.AcceptTime.After(b.AcceptTime) {
		return 1
	}
	return strings.Compare(a.ID, b.ID)
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		LimitBids:  tree.New[LimitOrder, Value](highToLowPrice),
		LimitAsks:  tree.New[LimitOrder, Value](lowToHighPrice),
		MarketBids: queue.New[MarketOrder](),
		MarketAsks: queue.New[MarketOrder](),
	}
}

type OrderBook struct {
	LimitBids  *tree.Tree[LimitOrder, Value]
	LimitAsks  *tree.Tree[LimitOrder, Value]
	MarketBids *queue.Queue[MarketOrder]
	MarketAsks *queue.Queue[MarketOrder]
}

func (o *OrderBook) GetMarketDepth() {

}

func (o *OrderBook) LimitAsk(id string, value Value, price Price) {
	newOrder := LimitOrder{
		ID:         id,
		Type:       OTAsk,
		Value:      value,
		Price:      price,
		AcceptTime: time.Now(),
	}
	// emit OrderAccepted(id, OTAsk, value, price, now);
	o.MarketBids.Iterate(func(order *MarketOrder, removeAndContinue func()) {
		if order.Value > value {
			order.Value -= value
			value = 0
			// emit OrderChanged(order.ID, order.Value)
			return
		}
		// order.Value <= value
		value -= order.Value
		// emit OrderExecuted(order.ID)
		removeAndContinue()
		return
	})
	if value == 0 {
		// emit OrderExecuted(id)
		return
	}
	o.LimitBids.Iterate(func(order LimitOrder, remainedValue *Value, removeAndContinue func()) {
		if order.Price < price {
			return
		}

		if *remainedValue > value {
			*remainedValue -= value
			value = 0
			// emit OrderChanged(order.ID, order.Value)
			return
		}
		// remainedValue <= value
		value -= *remainedValue
		// emit OrderExecuted(order.ID)
		removeAndContinue()
		return // get next
	})
	if value == 0 {
		// emit OrderExecuted(id)
		return
	}
	o.LimitAsks.Set(newOrder, value)
}

func (o *OrderBook) LimitBid(id string, value Value, price Price) {

}

func (o *OrderBook) MarketAsk(id string, value Value) {
	newOrder := MarketOrder{
		ID:         id,
		Type:       OTAsk,
		Value:      value,
		AcceptTime: time.Now(),
	}
	// emit OrderAccepted(id, OTAsk, value, price, now);
	o.MarketBids.Iterate(func(order *MarketOrder, removeAndContinue func()) {
		if order.Value > value {
			order.Value -= value
			value = 0
			// emit OrderChanged(order.ID, order.Value)
			return
		}
		// order.Value <= value
		value -= order.Value
		// emit OrderExecuted(order.ID)
		removeAndContinue()
		return
	})
	if value == 0 {
		// emit OrderExecuted(id)
		return
	}
	o.LimitBids.Iterate(func(order LimitOrder, remainedValue *Value, removeAndContinue func()) {
		if *remainedValue > value {
			*remainedValue -= value
			value = 0
			// emit OrderChanged(order.ID, order.Value)
			return
		}
		// remainedValue <= value
		value -= *remainedValue
		// emit OrderExecuted(order.ID)
		removeAndContinue()
		return // get next
	})
	if value == 0 {
		// emit OrderExecuted(id)
		return
	}
	o.MarketAsks.Add(newOrder)
}

func (o *OrderBook) MarketBid(id string, value Value) {

}
