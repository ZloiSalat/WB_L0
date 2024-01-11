package cache

import "WB/types"

// Cache ...
type Cache interface {
	Order() OrderCache
}

// OrderCache ...
type OrderCache interface {
	Load(orders map[string]*types.OrderJSON)
	Create(order *types.OrderJSON) error
	Find(orderUID string) (*types.OrderJSON, error)
}
