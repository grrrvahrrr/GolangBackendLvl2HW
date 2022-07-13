package kafka

import (
	"GoBeLvl2/recieveMsg/internal/entities"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/kafkapubsub"
)

type Kafka struct {
	ctx context.Context
	sub *pubsub.Subscription
}

func NewKafka(ctx context.Context) (*Kafka, error) {
	sub, err := pubsub.OpenSubscription(ctx, "kafka://MYG?topic=MY")
	if err != nil {
		return nil, err
	}
	return &Kafka{
		ctx: ctx,
		sub: sub,
	}, nil
}

func (k *Kafka) ListenToKafka(orderch chan entities.Order) {
	var order entities.Order
	for {
		msg, err := k.sub.Receive(k.ctx)
		if err != nil {
			log.Printf("Receiving message: %v", err)
			break
		}

		err = json.Unmarshal(msg.Body, &order)
		if err != nil {
			log.Printf("Json marshal error: %s", err)
			break
		}
		orderch <- order

		fmt.Printf("Received message: %s\n", string(msg.Body))
		fmt.Printf("Metadata: %v\n", msg.Metadata)

		msg.Ack()
	}
}

func (k *Kafka) SubClose() {
	k.sub.Shutdown(k.ctx)
}
