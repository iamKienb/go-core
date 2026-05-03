package kafkax

import "context"

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
