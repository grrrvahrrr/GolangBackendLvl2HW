package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nats-io/stan.go"
	"gocloud.dev/pubsub"

	_ "gocloud.dev/pubsub/rabbitpubsub"
	//_ "gocloud.dev/pubsub/kafkapubsub"
)

const (
	clusterID = "test-cluster"
	clientID  = "msgSender"
	model     = `{
		"order_uid": "b563feb7b2b84b6test3",
		"track_number": "WBILMTESTTRACK",
		"entry": "WBIL",
		"delivery": {
		  "name": "Test Testov",
		  "phone": "+9720000000",
		  "zip": "2639809",
		  "city": "Kiryat Mozkin",
		  "address": "Ploshad Mira 15",
		  "region": "Kraiot",
		  "email": "test@gmail.com"
		},
		"payment": {
		  "transaction": "b563feb7b2b84b6test3",
		  "request_id": "",
		  "currency": "USD",
		  "provider": "wbpay",
		  "amount": 1817,
		  "payment_dt": 1637907727,
		  "bank": "alpha",
		  "delivery_cost": 1500,
		  "goods_total": 317,
		  "custom_fee": 0
		},
		"items": [
		  {
			"chrt_id": 9934930,
			"track_number": "WBILMTESTTRACK",
			"price": 453,
			"rid": "ab4219087a764ae0btest",
			"name": "Mascaras",
			"sale": 30,
			"size": "0",
			"total_price": 317,
			"nm_id": 2389212,
			"brand": "Vivienne Sabo",
			"status": 202
		  }
		],
		"locale": "en",
		"internal_signature": "",
		"customer_id": "test",
		"delivery_service": "meest",
		"shardkey": "9",
		"sm_id": 99,
		"date_created": "2021-11-26T06:22:19Z",
		"oof_shard": "1"
	  }`
)

func main() {
	fmt.Println("Hello GoBeLvl2!")
	//useNats()

	ctx := context.Background()

	topic, err := WithRabbit(ctx)
	//topic, err := WithKafka(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer topic.Shutdown(ctx)

	err = topic.Send(ctx, &pubsub.Message{
		Body: []byte(model), // json

		Metadata: map[string]string{
			"language":   "en",
			"importance": "high",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func useNats() {
	sc, err := stan.Connect(
		clusterID,
		clientID,
		stan.NatsURL(stan.DefaultNatsURL),
	)
	if err != nil {
		log.Print(err)
	}

	defer sc.Close()

	sc.Publish("wbmodel", []byte(model))
}

func WithRabbit(ctx context.Context) (*pubsub.Topic, error) {
	// Under the hood - processing connect RABBIT_SERVER_URL=amqp://guest:guest@localhost:5672
	// UI for RabbitMQ - http://localhost:15672/
	// Another useful library - github.com/streadway/amqp
	return pubsub.OpenTopic(ctx, "rabbit://MY")
}

func WithKafka(ctx context.Context) (*pubsub.Topic, error) {
	// Under the hood - processing connect KAFKA_BROKERS=localhost:29092
	// UI for Kafka - http://localhost:8082/
	// Another useful library - github.com/Shopify/sarama
	return pubsub.OpenTopic(ctx, "kafka://MY")
}
