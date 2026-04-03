package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	topic  string
	retry  *RetryHandler
	dlq    *DLQ
	idem   *Idempotency
}

type Handler func(ctx context.Context, msg Message) error

func NewConsumer(brokers []string, topic string, groupID string, retry *RetryHandler, dlq *DLQ, idem *Idempotency) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        groupID,
			CommitInterval: 0,
			MinBytes:       10e3,
			MaxBytes:       10e6,
		}),
		topic: topic,
		retry: retry,
		dlq:   dlq,
		idem:  idem,
	}
}

func (c *Consumer) Consume(ctx context.Context, handler Handler) error {

	for {
		rawMsg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return err
			}
			continue
		}

		headers := make(map[string]string)

		for _, h := range rawMsg.Headers {
			headers[h.Key] = string(h.Value)
		}

		msg := Message{
			Key:       rawMsg.Key,
			Value:     rawMsg.Value,
			Headers:   headers,
			Timestamp: rawMsg.Time,
		}

		if d, ok := msg.Headers["x-delay"]; ok {
			if dur, err := time.ParseDuration(d); err == nil {
				time.Sleep(dur)
			}
		}

		if !c.idem.CheckAndSet(ctx, string(msg.Key)) {
			_ = c.reader.CommitMessages(ctx, rawMsg)
			continue
		}

		err = handler(ctx, msg)

		if err != nil {

			if c.retry != nil && c.retry.CanRetry(msg) {
				_ = c.retry.Handle(ctx, msg)
			} else if c.dlq != nil {
				_ = c.dlq.Send(ctx, msg)
			}

			_ = c.reader.CommitMessages(ctx, rawMsg)
			continue
		}

		if err := c.reader.CommitMessages(ctx, rawMsg); err != nil {
			log.Println("commit error:", err)
		}

	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
