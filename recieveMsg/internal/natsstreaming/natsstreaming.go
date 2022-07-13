package natsstreaming

import (
	"GoBeLvl2/recieveMsg/internal/entities"
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/stan.go"
)

const (
	clusterID = "test-cluster"
	clientID  = "event-store"
)

type NatsStreaming struct {
	sc  stan.Conn
	sub stan.Subscription
}

func NewNatsStreaming() (*NatsStreaming, error) {

	sc, err := stan.Connect(
		clusterID,
		clientID,
		stan.NatsURL(stan.DefaultNatsURL),
	)

	if err != nil {
		return nil, err
	}
	return &NatsStreaming{
		sc: sc,
	}, nil
}

func (ns *NatsStreaming) ListenToNats(ctx context.Context, orderch chan entities.Order) error {
	var order entities.Order

	sub, err := ns.sc.Subscribe("wbmodel",
		func(m *stan.Msg) {
			err := json.Unmarshal(m.Data, &order)
			if err != nil {
				log.Printf("Json marshal error: %s", err)
			}
			orderch <- order
		},
		stan.StartWithLastReceived())
	if err != nil {
		return err
	}

	ns.sub = sub

	return nil
}

func (ns *NatsStreaming) SubClose() {
	ns.sub.Close()
}
