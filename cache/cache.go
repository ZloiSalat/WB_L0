package cache

import (
	"WB"
	"WB/storage"
)

type Cache struct {
	orderCache *OrderCache
}

// New ...
func New(store storage.Storage) (*Cache, error) {
	orderCache, err := NewOrderCache(store)
	if err != nil {
		return nil, err
	}

	c := &Cache{
		orderCache: orderCache,
	}

	return c, nil
}

// Order ...
func (c *Cache) Order() OrderCache {
	if c.orderCache != nil {
		return c.orderCache
	}

	c.orderCache = &OrderCache{
		cache:  c,
		orders: make(map[string]*main.UserJSON),
	}

	return c.orderCache
}
