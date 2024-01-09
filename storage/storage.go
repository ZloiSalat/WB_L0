package storage

import (
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	db              *pgx.Conn
	orderRepository *OrderRepository
}

func NewS(db *pgx.Conn) *Storage {
	return &Storage{
		db: db,
	}
}

// Order ...
func (s *Storage) Order() OrderRepository {
	if s.orderRepository != nil {
		return s.orderRepository
	}

	s.orderRepository = &OrderRepository{
		store: s,
	}

	return s.orderRepository
}
