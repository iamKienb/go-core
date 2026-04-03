package kafka

import "context"

type DLQ struct {
	publisher *Publisher
}

func NewDLQ(pub *Publisher) *DLQ {
	return &DLQ{publisher: pub}
}

func (d *DLQ) Send(ctx context.Context, msg Message) error {
	if msg.Headers == nil {
		msg.Headers = make(map[string]string)
	}

	msg.Headers["x-dead"] = "true"

	return d.publisher.Publish(ctx, msg)
}
