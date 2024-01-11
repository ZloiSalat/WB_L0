package psqlstore

import (
	"WB/types"
	"context"
)

type OrderRepository struct {
	store *Store
}

func (s *OrderRepository) CreateUser(u *types.OrderJSON) error {
	_, err := s.store.db.Exec(context.Background(),
		"INSERT INTO orders (order_uid, data) "+
			"VALUES ($1, $2)",
		u.OrderUID, u.Data)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderRepository) FindAll() (map[string]*types.OrderJSON, error) {

	orders := make(map[string]*types.OrderJSON)

	query := "select order_uid, data from orders"

	rows, err := s.store.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		o := types.OrderJSON{}
		if err := rows.Scan(
			&o.OrderUID,
			&o.Data,
		); err != nil {
			return nil, err
		}
		orders[o.OrderUID] = &o
	}
	return orders, nil
}
