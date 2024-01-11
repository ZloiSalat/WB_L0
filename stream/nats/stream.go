package nats

import (
	"WB/cache"
	"WB/cache/mapcache"
	"WB/store"
	"WB/types"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
	"log"
)

type NatsStreaming struct {
	ns        stan.Conn
	store     store.Store
	validator *validator.Validate
	cache     cache.Cache
}

func NewNatsConnection(store store.Store, cache cache.Cache) (*NatsStreaming, error) {
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
		validator: validator.New(),
		cache:     cache,
	}, nil
}

func (sc *NatsStreaming) Subscribe() error {
	subject := "your-subject"
	_, err := sc.ns.Subscribe(subject, func(msg *stan.Msg) {
		//log.Printf("Message from stream %s", msg.Data)
		if err := msg.Ack(); err != nil {
			log.Println(err)
			return
		}
		var u types.User
		if err := json.Unmarshal(msg.Data, &u); err != nil {
			log.Println(err)
			return
		}

		if err := sc.validator.Struct(u); err != nil {
			log.Println(err)
			return
		}

		data, err := json.Marshal(u)
		if err != nil {
			log.Println(err)
			return
		}

		userCache := types.OrderJSON{
			OrderUID: u.OrderUid,
			Data:     data,
		}

		if err := sc.store.Order().CreateUser(&userCache); err != nil {
			log.Println(err)
			return
		}

		log.Printf("order with order_uid = %s stored to database\n", u.OrderUid)

		if err := sc.cache.Order().Create(&userCache); err != nil {
			log.Println(err)
			sc.cache, err = mapcache.New(sc.store)
			log.Printf("order with order_uid=%s stored to cache\n", userCache.OrderUID)
		}

	}, stan.SetManualAckMode())

	if err != nil {
		log.Fatalf("Error subscribing to NATS Streaming: %v", err)
	}

	return nil
}
