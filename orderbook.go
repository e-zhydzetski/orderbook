package orderbook

import (
	"time"

	"github.com/e-zhydzetski/orderbook/skiplist"

	"github.com/e-zhydzetski/orderbook/memtable"
)

// The result will be 0 if a == b, -1 if a < b, and +1 if a > b.
func lowToHighPrice(a PriceLimit, b PriceLimit) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// The result will be 0 if a == b, -1 if a < b, and +1 if a > b.
func highToLowPrice(a PriceLimit, b PriceLimit) int {
	if a > b {
		return -1
	}
	if a < b {
		return 1
	}
	return 0
}

func NewOrderBook() *OrderBook {
	orderGroupFactory := NewOrderGroupFactoryFunc()
	return &OrderBook{
		orderGroupFactory: orderGroupFactory,
		limitBids:         skiplist.New[PriceLimit, OrderGroup](10, highToLowPrice),
		// limitBids: tree.New[Order, Value](highToLowPrice),

		limitAsks: skiplist.New[PriceLimit, OrderGroup](10, lowToHighPrice),
		// limitAsks: tree.New[Order, Value](lowToHighPrice),

		marketBids: orderGroupFactory(),
		marketAsks: orderGroupFactory(),
		events:     NewEvents(100),
	}
}

type OrderBook struct {
	orderGroupFactory func() OrderGroup

	limitBids  memtable.MemTable[PriceLimit, OrderGroup]
	limitAsks  memtable.MemTable[PriceLimit, OrderGroup]
	marketBids OrderGroup
	marketAsks OrderGroup

	events *Events
}

//nolint:dupl,funlen // TODO refactor
func (o *OrderBook) Ask(id string, value Value, price PriceLimit) {
	now := time.Now()
	o.events.Emit(OrderAccepted{
		ID:         id,
		Type:       OTAsk,
		Value:      value,
		Price:      price,
		AcceptTime: now,
	})

	// match with market orders
	o.marketBids.Orders.Iterate(func(order *Order) memtable.IteratorAction {
		if value == 0 {
			return memtable.IAStop
		}
		if order.Value > value {
			order.Value -= value
			o.marketBids.TotalValue -= value
			value = 0
			o.events.Emit(OrderChanged{
				ID:    order.ID,
				Value: order.Value,
			})
			return memtable.IAStop
		}
		// order.Value <= value
		value -= order.Value
		o.marketBids.TotalValue -= order.Value
		o.events.Emit(OrderExecuted{
			ID: order.ID,
		})
		return memtable.IARemoveAndContinue
	})
	if value == 0 {
		o.events.Emit(OrderExecuted{
			ID: id,
		})
		return
	}

	// match with limit orders
	o.limitBids.Iterate(func(_ PriceLimit, orderGroup *OrderGroup) memtable.IteratorAction {
		cont := orderGroup.Orders.Iterate(func(order *Order) memtable.IteratorAction {
			if !price.IsMarket() {
				if order.Price < price {
					return memtable.IAStop
				}
			}
			if value == 0 {
				return memtable.IAStop
			}
			if order.Value > value {
				order.Value -= value
				orderGroup.TotalValue -= value
				value = 0
				o.events.Emit(OrderChanged{
					ID:    order.ID,
					Value: order.Value,
				})
				return memtable.IAStop
			}
			// remainedValue <= value
			value -= order.Value
			orderGroup.TotalValue -= order.Value
			o.events.Emit(OrderExecuted{
				ID: order.ID,
			})
			return memtable.IARemoveAndContinue
		})

		if cont {
			return memtable.IARemoveAndContinue
		}

		return memtable.IAStop
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
	newOrder := Order{
		ID:         id,
		Type:       OTAsk,
		Value:      value,
		Price:      price,
		AcceptTime: now,
	}
	if price.IsMarket() {
		o.marketAsks.Orders.PushHead(newOrder)
		o.marketAsks.TotalValue += newOrder.Value
	} else {
		o.limitAsks.Upsert(price, func() OrderGroup {
			og := o.orderGroupFactory()
			og.Orders.PushHead(newOrder)
			og.TotalValue += newOrder.Value
			return og
		}, func(og *OrderGroup) {
			og.Orders.PushHead(newOrder)
			og.TotalValue += newOrder.Value
		})
	}
	// o.events.PrintAll()
}

//nolint:dupl,funlen // TODO refactor
func (o *OrderBook) Bid(id string, value Value, price PriceLimit) {
	now := time.Now()
	o.events.Emit(OrderAccepted{
		ID:         id,
		Type:       OTBid,
		Value:      value,
		Price:      price,
		AcceptTime: now,
	})

	// match with market orders
	o.marketAsks.Orders.Iterate(func(order *Order) memtable.IteratorAction {
		if value == 0 {
			return memtable.IAStop
		}
		if order.Value > value {
			order.Value -= value
			o.marketAsks.TotalValue -= value
			value = 0
			o.events.Emit(OrderChanged{
				ID:    order.ID,
				Value: order.Value,
			})
			return memtable.IAStop
		}
		// order.Value <= value
		value -= order.Value
		o.marketAsks.TotalValue -= order.Value
		o.events.Emit(OrderExecuted{
			ID: order.ID,
		})
		return memtable.IARemoveAndContinue
	})
	if value == 0 {
		o.events.Emit(OrderExecuted{
			ID: id,
		})
		return
	}

	// match with limit orders
	o.limitAsks.Iterate(func(_ PriceLimit, orderGroup *OrderGroup) memtable.IteratorAction {
		cont := orderGroup.Orders.Iterate(func(order *Order) memtable.IteratorAction {
			if !price.IsMarket() {
				if order.Price > price {
					return memtable.IAStop
				}
			}
			if value == 0 {
				return memtable.IAStop
			}
			if order.Value > value {
				order.Value -= value
				orderGroup.TotalValue -= value
				value = 0
				o.events.Emit(OrderChanged{
					ID:    order.ID,
					Value: order.Value,
				})
				return memtable.IAStop
			}
			// remainedValue <= value
			value -= order.Value
			orderGroup.TotalValue -= order.Value
			o.events.Emit(OrderExecuted{
				ID: order.ID,
			})
			return memtable.IARemoveAndContinue
		})

		if cont {
			return memtable.IARemoveAndContinue
		}

		return memtable.IAStop
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
	newOrder := Order{
		ID:         id,
		Type:       OTBid,
		Value:      value,
		Price:      price,
		AcceptTime: now,
	}
	if price.IsMarket() {
		o.marketBids.Orders.PushHead(newOrder)
		o.marketBids.TotalValue += newOrder.Value
	} else {
		o.limitBids.Upsert(price, func() OrderGroup {
			og := o.orderGroupFactory()
			og.Orders.PushHead(newOrder)
			og.TotalValue += newOrder.Value
			return og
		}, func(og *OrderGroup) {
			og.Orders.PushHead(newOrder)
			og.TotalValue += newOrder.Value
		})
	}
	// o.events.PrintAll()
}
