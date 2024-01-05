package main

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Storage interface {
	CreateUser(*User) error
	DeleteSegment(string) error
	GetActiveUsers() (*User, error)
}

type PostgresStore struct {
	db *pgx.Conn
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "postgres://wb_user:wb_password@localhost:5434/wb_db"
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	//defer conn.Close(context.Background())
	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: conn,
	}, nil
}

func (s *PostgresStore) DeleteSegment(slug string) error {
	query := "delete from segments where segment_name=$1"
	_, err := s.db.Query(context.Background(), query, slug)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) CreateUser(u *User) error {
	_, err := s.db.Exec(context.Background(),
		"INSERT INTO orders (order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)",
		u.OrderUid, u.TrackNumber, u.Entry, u.Delivery, u.Payment, u.Items, u.Locale, u.InternalSignature, u.CustomerID, u.DeliveryService, u.Shardkey, u.SmID, u.DateCreated, u.OofShard)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) GetActiveUsers() (*User, error) {
	//TODO implement me
	panic("implement me")
}

/*func (s *PostgresStore) GetActiveSegments(id int) (*User, error) {
	query := `select s.segment_name
			  from segments s
			  inner join user_segments us on s.id = us.segment_id
			  where us.user_id = $1;`
	rows, err := s.db.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}

	var concatenatedSegment string // Строка для конкатенации сегментов.

	for rows.Next() {
		var segmentName string
		if err := rows.Scan(&segmentName); err != nil {
			return nil, err
		}
		concatenatedSegment += segmentName + "," // Конкатенируем сегменты с запятой в качестве разделителя.
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Удалите последнюю запятую, если она есть.
	if len(concatenatedSegment) > 0 && concatenatedSegment[len(concatenatedSegment)-1] == ',' {
		concatenatedSegment = concatenatedSegment[:len(concatenatedSegment)-1]
	}

	return &User{
		Segment: concatenatedSegment,
	}, nil

}*/
