package orderbook

import (
	"strings"
	"time"

	"github.com/e-zhydzetski/strips-tt/orderbook/queue"
	"github.com/e-zhydzetski/strips-tt/orderbook/tree"
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
		limitBids:  tree.New[LimitOrder, Value](highToLowPrice),
		limitAsks:  tree.New[LimitOrder, Value](lowToHighPrice),
		marketBids: queue.New[MarketOrder](),
		marketAsks: queue.New[MarketOrder](),
		events:     NewEvents(100),
	}
}

type OrderBook struct {
	limitBids  *tree.Tree[LimitOrder, Value] // TODO tree of queues
	limitAsks  *tree.Tree[LimitOrder, Value] // TODO tree of queues
	marketBids *queue.Queue[MarketOrder]
	marketAsks *queue.Queue[MarketOrder]
	events     *Events
}

//nolint:dupl // TODO refactor
func (o *OrderBook) Ask(id string, value Value, price PriceLimit) {
	now := time.Now()
	o.events.Emit(OrderAccepted{
		ID:         id,
		Type:       OTAsk,
		Value:      value,
		Price:      price,
		AcceptTime: now,
	})
	o.marketBids.Iterate(func(order *MarketOrder) queue.IteratorAction {
		if order.Value > value {
			order.Value -= value
			value = 0
			o.events.Emit(OrderChanged{
				ID:    order.ID,
				Value: order.Value,
			})
			return queue.IAStop
		}
		// order.Value <= value
		value -= order.Value
		o.events.Emit(OrderExecuted{
			ID: order.ID,
		})
		return queue.IARemoveAndContinue
	})
	if value == 0 {
		o.events.Emit(OrderExecuted{
			ID: id,
		})
		return
	}
	o.limitBids.Iterate(func(order LimitOrder, remainedValue *Value) tree.IteratorAction {
		if !price.IsMarket() {
			if order.Price < price {
				return tree.IAStop
			}
		}

		if *remainedValue > value {
			*remainedValue -= value
			value = 0
			o.events.Emit(OrderChanged{
				ID:    order.ID,
				Value: order.Value,
			})
			return tree.IAStop
		}
		// remainedValue <= value
		value -= *remainedValue
		o.events.Emit(OrderExecuted{
			ID: order.ID,
		})
		return tree.IARemoveAndContinue
	})
	if value == 0 {
		o.events.Emit(OrderExecuted{
			ID: id,
		})
		return
	}
	o.events.Emit(OrderChanged{
		ID:    id,
		Value: value,
	})

	if price.IsMarket() {
		newOrder := MarketOrder{
			ID:         id,
			Type:       OTAsk,
			Value:      value,
			AcceptTime: now,
		}
		o.marketAsks.Add(newOrder)
	} else {
		newOrder := LimitOrder{
			ID:         id,
			Type:       OTAsk,
			Value:      value,
			Price:      price,
			AcceptTime: now,
		}
		o.limitAsks.Set(newOrder, value)
	}
	// o.events.PrintAll()
}

//nolint:dupl // TODO refactor
func (o *OrderBook) Bid(id string, value Value, price PriceLimit) {
	now := time.Now()
	o.events.Emit(OrderAccepted{
		ID:         id,
		Type:       OTBid,
		Value:      value,
		Price:      price,
		AcceptTime: now,
	})
	o.marketAsks.Iterate(func(order *MarketOrder) queue.IteratorAction {
		if order.Value > value {
			order.Value -= value
			value = 0
			o.events.Emit(OrderChanged{
				ID:    order.ID,
				Value: order.Value,
			})
			return queue.IAStop
		}
		// order.Value <= value
		value -= order.Value
		o.events.Emit(OrderExecuted{
			ID: order.ID,
		})
		return queue.IARemoveAndContinue
	})
	if value == 0 {
		o.events.Emit(OrderExecuted{
			ID: id,
		})
		return
	}
	o.limitAsks.Iterate(func(order LimitOrder, remainedValue *Value) tree.IteratorAction {
		if !price.IsMarket() {
			if order.Price < price {
				return tree.IAStop
			}
		}

		if *remainedValue > value {
			*remainedValue -= value
			value = 0
			o.events.Emit(OrderChanged{
				ID:    order.ID,
				Value: order.Value,
			})
			return tree.IAStop
		}
		// remainedValue <= value
		value -= *remainedValue
		o.events.Emit(OrderExecuted{
			ID: order.ID,
		})
		return tree.IARemoveAndContinue
	})
	if value == 0 {
		o.events.Emit(OrderExecuted{
			ID: id,
		})
		return
	}
	o.events.Emit(OrderChanged{
		ID:    id,
		Value: value,
	})

	if price.IsMarket() {
		newOrder := MarketOrder{
			ID:         id,
			Type:       OTBid,
			Value:      value,
			AcceptTime: now,
		}
		o.marketBids.Add(newOrder)
	} else {
		newOrder := LimitOrder{
			ID:         id,
			Type:       OTBid,
			Value:      value,
			Price:      price,
			AcceptTime: now,
		}
		o.limitBids.Set(newOrder, value)
	}
	// o.events.PrintAll()
}
