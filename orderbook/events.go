package orderbook

import (
	"log"
	"time"
)

func NewEvents(size int) *Events {
	if size == 0 {
		panic("event buffer with zero size")
	}
	return &Events{
		buff:    make([]any, size),
		nextIdx: 0,
	}
}

type Events struct {
	buff    []any
	nextIdx int
}

func (e *Events) Emit(event any) {
	e.buff[e.nextIdx] = event
	e.nextIdx = (e.nextIdx + 1) % len(e.buff)
}

func (e *Events) PrintAll() {
	first := true
	for i := e.nextIdx; first || i != e.nextIdx; i = (i + 1) % len(e.buff) {
		first = false
		if e.buff[i] == nil {
			continue
		}
		log.Printf("%#v", e.buff[i])
	}
}

type OrderAccepted struct {
	ID         string
	Type       OrderType
	Value      Value
	Price      PriceLimit
	AcceptTime time.Time
}

type OrderExecuted struct {
	ID string
}

type OrderChanged struct {
	ID    string
	Value Value
}
