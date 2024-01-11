package psqlstore

import (
	"WB/store"
	"context"
	"github.com/jackc/pgx/v5"
)

type Store struct {
	db              *pgx.Conn
	orderRepository *OrderRepository
}

func NewPostgresStore() (*Store, error) {
	connStr := "postgres://wb_user:wb_password@localhost:5434/wb_db"
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	//defer conn.Close(context.Background())

	if err != nil {
		return nil, err
	}

	return &Store{
		db: conn,
	}, nil
}

func (s *Store) Order() store.OrderRepository {
	if s.orderRepository != nil {
		return s.orderRepository
	}

	s.orderRepository = &OrderRepository{
		store: s,
	}

	return s.orderRepository
}
