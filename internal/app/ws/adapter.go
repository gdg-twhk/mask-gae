package ws

import (
	"context"
	"errors"
	"sync"

	broker "github.com/nats-io/nats.go"
)

var (
	// ErrFailedMessagePublish indicates that message publishing failed.
	ErrFailedMessagePublish = errors.New("failed to publish message")

	// ErrFailedSubscription indicates that client couldn't subscribe to specified channel.
	ErrFailedSubscription = errors.New("failed to subscribe to a channel")

	// ErrFailedConnection indicates that service couldn't connect to message broker.
	ErrFailedConnection = errors.New("failed to connect to message broker")
)

// Service specifies web socket service API.
type Service interface {
	Publish(ctx context.Context, topic string, msg string) error

	// Subscribes to channel with specified id.
	Subscribe(topic string, channel *Channel) error
}

// Channel is used for receiving and sending messages.
type Channel struct {
	Messages chan string
	Closed   chan bool
	closed   bool
	mutex    sync.Mutex
}

// NewChannel instantiates empty channel.
func NewChannel() *Channel {
	return &Channel{
		Messages: make(chan string),
		Closed:   make(chan bool),
		closed:   false,
		mutex:    sync.Mutex{},
	}
}

// Send method send message over Messages channel.
func (channel *Channel) Send(msg string) {
	channel.mutex.Lock()
	defer channel.mutex.Unlock()

	if !channel.closed {
		channel.Messages <- msg
	}
}

// Close channel and stop message transfer.
func (channel *Channel) Close() {
	channel.mutex.Lock()
	defer channel.mutex.Unlock()

	channel.closed = true
	channel.Closed <- true
	close(channel.Messages)
	close(channel.Closed)
}

var _ Service = (*adapterService)(nil)

type adapterService struct {
	pubsub Service
}

// New instantiates the WS adapter implementation.
func New(pubsub Service) Service {
	return &adapterService{pubsub: pubsub}
}

func (as *adapterService) Publish(ctx context.Context, topic string, msg string) error {
	if err := as.pubsub.Publish(ctx, topic, msg); err != nil {
		switch err {
		case broker.ErrConnectionClosed, broker.ErrInvalidConnection:
			return ErrFailedConnection
		default:
			return ErrFailedMessagePublish
		}
	}
	return nil
}

func (as *adapterService) Subscribe(topic string, channel *Channel) error {
	if err := as.pubsub.Subscribe(topic, channel); err != nil {
		return ErrFailedSubscription
	}

	return nil
}
