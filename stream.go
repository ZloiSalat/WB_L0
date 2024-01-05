package main

import (
	"github.com/nats-io/stan.go"
	"log"
)

var exitSignal = make(chan struct{})

type Stream interface {
	subscribe() error
}

type NatsStreaming struct {
	ns   stan.Conn
	Data []byte `json:"data"`
}

func NewNatsConnection() (*NatsStreaming, error) {
	clientID := "your-client-id--"
	clusterID := "my-cluster"
	natsURL := "nats://localhost:4222" // Update with your NATS Streaming server URL

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL), stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
		log.Printf("Connection lost, reason: %v", reason)
	}))
	if err != nil {
		log.Fatalf("Error connecting to NATS Streaming: %v", err)
	}

	//defer sc.Close()
	log.Println("Connected to nats")
	return &NatsStreaming{
		ns:   sc,
		Data: make([]byte, 0),
	}, nil
}

func (sc *NatsStreaming) subscribe() error {
	log.Print("Subsribing!")
	subject := "your-subject"
	sub, err := sc.ns.Subscribe(subject, func(msg *stan.Msg) {
		//log.Printf("Message from stream %s", msg.Data)
		sc.Data = append(sc.Data, msg.Data...)
		log.Printf("Message from stream %s", sc.Data)

	})
	if err != nil {
		log.Fatalf("Error subscribing to NATS Streaming: %v", err)
	}
	defer sub.Close()
	if err != nil {
		log.Fatalf("Error closing subscription to NATS Streaming: %v", err)
	}

	<-exitSignal

	return nil
}
