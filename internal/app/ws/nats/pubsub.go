package nats

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	broker "github.com/nats-io/nats.go"
	"github.com/sony/gobreaker"

	"github.com/cage1016/mask/internal/app/ws"
)

const (
	maxFailedReqs   = 3
	maxFailureRatio = 0.6
)

var _ ws.Service = (*natsPubSub)(nil)

type natsPubSub struct {
	nc     *broker.Conn
	cb     *gobreaker.CircuitBreaker
	logger log.Logger
}

// New instantiates NATS message publisher.
func New(nc *broker.Conn, logger log.Logger) ws.Service {
	st := gobreaker.Settings{
		Name: "NATS",
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			fr := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= maxFailedReqs && fr >= maxFailureRatio
		},
	}
	cb := gobreaker.NewCircuitBreaker(st)
	return &natsPubSub{
		nc:     nc,
		cb:     cb,
		logger: logger,
	}
}

func (pubsub *natsPubSub) Publish(_ context.Context, topic string, msg string) error {
	return pubsub.nc.Publish(topic, []byte(msg))
}

func (pubsub *natsPubSub) Subscribe(topic string, channel *ws.Channel) error {
	var sub *broker.Subscription

	pubsub.nc.Subscribe(topic, func(msg *broker.Msg) {
		if msg == nil {
			level.Warn(pubsub.logger).Log("msg", "Received nil message")
			return
		}

		// Sends message to messages channel
		channel.Send(string(msg.Data))
	})
	pubsub.nc.Flush()

	// Check if subscription should be closed
	go func() {
		<-channel.Closed
		sub.Unsubscribe()
	}()

	err := pubsub.nc.LastError()
	if err != nil {
		level.Warn(pubsub.logger).Log("lastError", err)
		return err
	}

	return nil
}
