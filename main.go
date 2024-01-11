package main

import (
	"WB/app"
	"WB/cache/mapcache"
	store2 "WB/store/psqlstore"
	stream2 "WB/stream/nats"
	"log"
)

func main() {

	store, err := store2.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	cache, err := mapcache.New(store)
	if err != nil {
		log.Fatal(err)
	}
	stream, err := stream2.NewNatsConnection(store, cache)
	if err != nil {
		log.Fatal(err)
	}
	if err := stream.Subscribe(); err != nil {
		log.Fatal(err)
	}

	server := app.NewAPIServer(":3000", store, cache)
	server.Run()

}
