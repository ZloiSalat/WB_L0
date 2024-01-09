package cache

import (
	"WB"
	"WB/storage"
	"errors"
)

// OrderCache ...
type OrderCache struct {
	orders map[string]*main.UserJSON
	cache  *Cache
}

// NewOrderCache ...
func NewOrderCache(store storage.Storage) (*OrderCache, error) {
	orders, err := store.FindAll()
	if err != nil {
		return nil, err
	}

	c := &OrderCache{
		orders: orders,
	}

	return c, nil
}

func (c *OrderCache) Load(orders map[string]*main.UserJSON) {
	c.orders = orders
}

func (c *OrderCache) Create(order *main.UserJSON) error {
	if _, ok := c.orders[order.OrderUID]; ok == true {
		return errors.New("already exists")
	}

	c.orders[order.OrderUID] = order
	return nil
}

func (c *OrderCache) Find(orderUID string) (*main.UserJSON, error) {
	if order, ok := c.orders[orderUID]; ok == true {
		return order, nil
	}

	return nil, errors.New("record not found")
}
