package main

import (
	"log"
)

func main() {

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	stream, err := NewNatsConnection()
	if err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":3000", store, stream)

	server.Run()

	close(exitSignal)
}
