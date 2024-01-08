package main

import (
	"log"
)

func main() {

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":3000", store)
	stream, err := NewNatsConnection(store, server)
	if err != nil {
		log.Fatal(err)
	}

	if err := stream.Subscribe(); err != nil {
		log.Fatal(err)
	}

	server.Run()

}
