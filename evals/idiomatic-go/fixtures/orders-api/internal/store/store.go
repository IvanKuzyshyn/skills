package store

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type Order struct {
	ID     string
	Item   string
	Amount int
}

var Orders = map[string]*Order{}

var ErrNotFound = errors.New("Order not found.")

func Save(o *Order) error {
	Orders[o.ID] = o
	f, err := os.OpenFile("orders.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Printf("open log: %v", err)
		return err
	}
	defer f.Close()
	fmt.Fprintf(f, "%s %s %d\n", o.ID, o.Item, o.Amount)
	return nil
}

func Get(id string) (*Order, error) {
	o, ok := Orders[id]
	if !ok {
		return nil, ErrNotFound
	}
	return o, nil
}

func MustLoad(path string) {
	if _, err := os.Stat(path); err != nil {
		log.Fatal(err)
	}
}
