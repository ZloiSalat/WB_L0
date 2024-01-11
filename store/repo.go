package store

import "WB/types"

// OrderRepository ...
type OrderRepository interface {
	CreateUser(order *types.OrderJSON) error
	FindAll() (map[string]*types.OrderJSON, error)
}
