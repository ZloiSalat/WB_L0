package mapcache

import (
	"WB/store"
	"WB/types"
	"errors"
	"log"
)

// OrderCache ...
type OrderCache struct {
	orders map[string]*types.OrderJSON
	cache  *Cache
}

// NewOrderCache ...
func NewOrderCache(store store.Store) (*OrderCache, error) {
	orders, err := store.Order().FindAll()
	if err != nil {
		return nil, err
	}

	c := &OrderCache{
		orders: orders,
	}

	log.Println("Initialized OrderCache with orders:", orders)

	return c, nil
}

func (c *OrderCache) Load(orders map[string]*types.OrderJSON) {
	c.orders = orders
}

func (c *OrderCache) Create(order *types.OrderJSON) error {
	if _, ok := c.orders[order.OrderUID]; ok == true {
		return errors.New("already exists")
	}

	c.orders[order.OrderUID] = order
	return nil
}

func (c *OrderCache) Find(orderUID string) (*types.OrderJSON, error) {
	if order, ok := c.orders[orderUID]; ok == true {
		return order, nil
	}

	return nil, errors.New("record not found")
}
