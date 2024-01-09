package storage

import (
	"WB"
	"context"
)

/*type Storage interface {
	CreateUser(UserJSON) error
	DeleteSegment(string) error
	GetActiveUsers() (*User, error)
	FindAll() (map[string]*UserJSON, error)
}*/

type OrderRepository struct {
	store *Storage
}

func (s *OrderRepository) CreateUser(u *main.UserJSON) error {
	_, err := s.store.db.Exec(context.Background(),
		"INSERT INTO orders (order_uid, data) "+
			"VALUES ($1, $2)",
		u.OrderUID, u.Data)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderRepository) FindAll() (map[string]*main.UserJSON, error) {

	orders := make(map[string]*main.UserJSON)

	query := "select order_uid, data from orders"

	rows, err := s.store.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := main.UserJSON{}
		if err := rows.Scan(
			&u.OrderUID,
			&u.Data,
		); err != nil {
			return nil, err
		}
		orders[u.OrderUID] = &u
	}
	return orders, nil
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
