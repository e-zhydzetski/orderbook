package orderbook

import (
	"log"
	"testing"
)

func TestEvents(t *testing.T) {
	e := NewEvents(3)
	e.PrintAll() // nothing
	log.Println("----")
	e.Emit(1)
	e.PrintAll() // 1
	log.Println("----")
	e.Emit(2)
	e.Emit(3)
	e.PrintAll() // 1, 2, 3
	log.Println("----")
	e.Emit(4)
	e.PrintAll() // 2, 3, 4
}
