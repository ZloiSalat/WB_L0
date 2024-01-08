package main

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
	"log"
)

type Stream interface {
	Subscribe() error
	stopSubscription()
}

type NatsStreaming struct {
	ns        stan.Conn
	store     Storage
	server    *APIServer
	validator *validator.Validate
}

func NewNatsConnection(store Storage, server *APIServer) (*NatsStreaming, error) {
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
		ns:        sc,
		store:     store,
		server:    server,
		validator: validator.New(),
	}, nil
}

func (sc *NatsStreaming) Subscribe() error {
	subject := "your-subject"
	_, err := sc.ns.Subscribe(subject, func(msg *stan.Msg) {
		log.Printf("Message from stream %s", msg.Data)
		if err := msg.Ack(); err != nil {
			log.Println(err)
			return
		}
		var u User
		if err := json.Unmarshal(msg.Data, &u); err != nil {
			log.Println(err)
			return
		}

		log.Println(u)

		if err := sc.validator.Struct(u); err != nil {
			log.Println(err)
			return
		}

		userJSON, _ := NewUser(
			u.OrderUid,
			u.TrackNumber,
			u.Entry,
			u.Delivery,
			u.Payment,
			u.Items,
			u.Locale,
			u.InternalSignature,
			u.CustomerID,
			u.DeliveryService,
			u.Shardkey,
			u.SmID,
			u.DateCreated,
			u.OofShard,
		)

		if err := sc.store.CreateUser(userJSON); err != nil {
			log.Println(err)
			return
		}

		log.Printf("order with order_uid = %s stored to database\n", u.OrderUid)

	}, stan.SetManualAckMode())

	if err != nil {
		log.Fatalf("Error subscribing to NATS Streaming: %v", err)
	}

	return nil
}
