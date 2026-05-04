package kafkax

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaXService interface {
	Ping(ctx context.Context) error
	Dialer() *kafka.Dialer
	Transport() *kafka.Transport
	Brokers() []string
	ClientID() string
	Close() error
}

type ProducerPublisher interface {
	Publish(ctx context.Context, msg Message) error
	PublishBatch(ctx context.Context, messages []Message) error
}

type ConsumerHandler interface {
	Handle(ctx context.Context, msg Message) error
}

type ConsumerHandlerFunc func(ctx context.Context, msg Message) error

func (f ConsumerHandlerFunc) Handle(ctx context.Context, msg Message) error {
	return f(ctx, msg)
}
