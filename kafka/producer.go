package kafkax

import (
	"context"
	"errors"
	"fmt"
	"strings"

	configx "github.com/iamKienb/shopify-go-platform/config"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer(client *KafkaX, cfg configx.ProducerConfig) (*Producer, error) {
	if client == nil {
		return nil, errors.New("kafka producer: client must not be nil")
	}
	if strings.TrimSpace(cfg.Topic) == "" {
		return nil, errors.New("kafka producer: topic must not be empty")
	}

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(client.Brokers()...),
		Topic:                  cfg.Topic,
		Balancer:               &kafka.LeastBytes{},
		BatchTimeout:           cfg.BatchTimeout,
		BatchBytes:             cfg.BatchBytes,
		RequiredAcks:           cfg.RequiredAcks,
		AllowAutoTopicCreation: cfg.AllowAutoTopic,
		Async:                  false,
		MaxAttempts:            cfg.MaxAttempts,
		Transport:              client.Transport(),
		ReadTimeout:            cfg.ReadTimeout,
		WriteTimeout:           cfg.WriteTimeout,
		Compression:            resolveCompression(cfg.Compression),
	}

	return &Producer{
		writer: writer,
		topic:  cfg.Topic,
	}, nil
}

func (p *Producer) Publish(ctx context.Context, msg Message) error {
	return p.PublishBatch(ctx, []Message{msg})
}

func (p *Producer) PublishBatch(ctx context.Context, messages []Message) error {
	if len(messages) == 0 {
		return nil
	}

	kafkaMessages := make([]kafka.Message, 0, len(messages))
	for _, msg := range messages {
		if msg.Topic == "" {
			msg.Topic = p.topic
		}
		kafkaMessages = append(kafkaMessages, toKafkaMessage(msg))
	}

	if err := p.writer.WriteMessages(ctx, kafkaMessages...); err != nil {
		return fmt.Errorf("kafka publish failed: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	if p.writer == nil {
		return nil
	}
	return p.writer.Close()
}

func resolveCompression(name string) compress.Compression {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "gzip":
		return compress.Gzip
	case "snappy":
		return compress.Snappy
	case "lz4":
		return compress.Lz4
	case "zstd":
		return compress.Zstd
	default:
		return compress.None
	}
}
