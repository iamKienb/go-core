package kafka

import (
	"context"
	"time"

	"github.com/iamKienb/shopify-go-platform/utils"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Publisher struct {
	writer *kafka.Writer
	topic  string
}

func NewPublisher(brokers []string, topic string) *Publisher {

	return &Publisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			RequiredAcks: kafka.RequireAll,
			MaxAttempts:  5,
			BatchSize:    100,
			BatchTimeout: 10 * time.Millisecond,
			Balancer:     &kafka.LeastBytes{},
		},
		topic: topic,
	}

}

func (p *Publisher) Publish(ctx context.Context, msg Message) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	headers := []kafka.Header{}

	for k, v := range msg.Headers {
		headers = append(headers, kafka.Header{
			Key:   k,
			Value: []byte(v),
		})
	}

	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:     msg.Key,
		Value:   msg.Value,
		Headers: headers,
		Time:    msg.Timestamp,
	})

	if err != nil {
		utils.GetLogger().Error(
			"Kafka publish failed",
			zap.String("topic", p.topic),
			zap.ByteString("key", msg.Key),
			zap.Any("error", err),
		)
	}

	return err
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}
