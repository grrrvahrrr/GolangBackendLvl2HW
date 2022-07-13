package rabbit

import (
	"GoBeLvl2/recieveMsg/internal/entities"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/rabbitpubsub"
)

type Rabbit struct {
	ctx context.Context
	sub *pubsub.Subscription
}

func NewRabbit(ctx context.Context) (*Rabbit, error) {
	sub, err := pubsub.OpenSubscription(ctx, "rabbit://MY")
	if err != nil {
		return nil, err
	}
	return &Rabbit{
		ctx: ctx,
		sub: sub,
	}, nil
}

func (r *Rabbit) ListenToRabbit(orderch chan entities.Order) {
	var order entities.Order
	for {
		msg, err := r.sub.Receive(r.ctx)
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

func (r *Rabbit) SubClose() {
	r.sub.Shutdown(r.ctx)
}
