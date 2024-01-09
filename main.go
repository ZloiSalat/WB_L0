package main

import (
	cache2 "WB/cache"
	"WB/storage"
	"context"
	"github.com/jackc/pgx/v5"
	"log"
)

func main() {

	db, err := db()
	if err != nil {
		log.Fatal(err)
	}
	store := storage.NewS(db)
	cache, err := cache2.New(*store)
	if err != nil {
		log.Fatal(err)
	}
	stream, err := NewNatsConnection(store, cache)
	if err != nil {
		log.Fatal(err)
	}
	if err := stream.Subscribe(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":3000", store)
	server.Run()

}

func db() (*pgx.Conn, error) {
	connStr := "postgres://wb_user:wb_password@localhost:5434/wb_db"
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())
	return conn, nil
}
