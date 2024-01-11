package mapcache

import (
	"WB/cache"
	"WB/store"
	"WB/types"
)

type Cache struct {
	orderCache *OrderCache
}

// New ...
func New(store store.Store) (*Cache, error) {
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
func (c *Cache) Order() cache.OrderCache {
	if c.orderCache != nil {
		return c.orderCache
	}

	c.orderCache = &OrderCache{
		cache:  c,
		orders: make(map[string]*types.OrderJSON),
	}

	return c.orderCache
}
