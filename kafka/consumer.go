package kafkax

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type ConsumerConfig struct {
	GroupID      string
	Topic        string
	DLQTopic     string
	MinBytes     int
	MaxBytes     int
	MaxWait      time.Duration
	MaxAttempts  int
	RetryBackoff time.Duration
	Logger       *slog.Logger
}

type Consumer struct {
	reader      *kafka.Reader
	handler     ConsumerHandler
	dlqProducer *Producer
	cfg         ConsumerConfig
}

func NewConsumer(client *KafkaX, cfg ConsumerConfig, handler ConsumerHandler) (*Consumer, error) {
	if client == nil || handler == nil {
		return nil, errors.New("kafka consumer: client and handler must not be nil")
	}

	if strings.TrimSpace(cfg.Topic) == "" || strings.TrimSpace(cfg.GroupID) == "" {
		return nil, errors.New("kafka consumer: topic and group id must not be empty")
	}

	cfg = normalizeConsumerConfig(cfg)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  client.Brokers(),
		GroupID:  cfg.GroupID,
		Topic:    cfg.Topic,
		MinBytes: cfg.MinBytes,
		MaxBytes: cfg.MaxBytes,
		MaxWait:  cfg.MaxWait,
		Dialer:   client.Dialer(),
	})

	var dlqProducer *Producer
	if cfg.DLQTopic != "" {
		var err error
		dlqProducer, err = NewProducer(client, ProducerConfig{
			Topic:        cfg.DLQTopic,
			Balancer:     "least_bytes",
			BatchTimeout: 50 * time.Millisecond,
			MaxAttempts:  3,
		})
		if err != nil {
			return nil, fmt.Errorf("kafka consumer dlq producer init failed: %w", err)
		}
	}

	return &Consumer{
		reader:      reader,
		handler:     handler,
		dlqProducer: dlqProducer,
		cfg:         cfg,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	c.logInfo(ctx, "kafka consumer started", slog.String("topic", c.cfg.Topic))
	for {
		if ctx.Err() != nil {
			return nil
		}

		if err := c.consumeOne(ctx); err != nil {
			if ctx.Err() != nil {
				return nil
			}

			c.logError(ctx, "kafka consume failed", slog.String("error", err.Error()))
			time.Sleep(c.cfg.RetryBackoff)
		}
	}
}

func (c *Consumer) consumeOne(ctx context.Context) error {
	kmsg, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return fmt.Errorf("kafka fetch failed: %w", err)
	}

	msg := fromKafkaMessage(c.cfg.Topic, kmsg)

	if err := c.handleWithRetry(ctx, msg); err != nil {
		if c.dlqProducer != nil {
			if dlqErr := c.publishDLQ(ctx, msg, err); dlqErr != nil {
				return fmt.Errorf("failed to send to DLQ: %w (original error: %v)", dlqErr, err)
			}
		} else {
			return fmt.Errorf("handler failed and no DLQ configured: %w", err)
		}

	}

	if err := c.reader.CommitMessages(ctx, kmsg); err != nil {
		return fmt.Errorf("kafka commit failed: %w", err)
	}

	return nil
}

func (c *Consumer) handleWithRetry(ctx context.Context, msg Message) error {
	var lastErr error
	for attempt := 1; attempt <= c.cfg.MaxAttempts; attempt++ {
		if attempt > 1 {
			time.Sleep(time.Duration(attempt-1) * c.cfg.RetryBackoff)
		}

		msg.SetHeader(HeaderRetryCount, strconv.Itoa(attempt-1))
		if err := c.handler.Handle(ctx, msg); err != nil {
			lastErr = err
			c.logWarn(ctx, "kafka handler attempt failed",
				slog.String("topic", c.cfg.Topic),
				slog.String("group_id", c.cfg.GroupID),
				slog.Int("attempt", attempt),
				slog.String("error", err.Error()),
			)
			continue
		}

		return nil
	}

	return fmt.Errorf("kafka handler failed after %d attempts: %w", c.cfg.MaxAttempts, lastErr)
}

func (c *Consumer) publishDLQ(ctx context.Context, msg Message, cause error) error {
	if c.dlqProducer == nil {
		return cause
	}

	dlqMsg := msg
	dlqMsg.Topic = c.cfg.DLQTopic
	dlqMsg.SetHeader(HeaderOriginalTopic, c.cfg.Topic)
	dlqMsg.SetHeader(HeaderOriginalOffset, strconv.FormatInt(msg.Offset, 10))
	dlqMsg.SetHeader(HeaderOriginalGroupID, c.cfg.GroupID)
	dlqMsg.SetHeader(HeaderFailureReason, cause.Error())
	dlqMsg.SetHeader(HeaderFailureAt, time.Now().UTC().Format(time.RFC3339))

	return c.dlqProducer.Publish(ctx, dlqMsg)
}

func (c *Consumer) Close() error {
	var result error
	if c.reader != nil {
		result = c.reader.Close()
	}
	if c.dlqProducer != nil {
		if err := c.dlqProducer.Close(); err != nil && result == nil {
			result = err
		}
	}
	return result
}

func (c *Consumer) logInfo(ctx context.Context, msg string, args ...any) {
	if c.cfg.Logger != nil {
		c.cfg.Logger.InfoContext(ctx, msg, args...)
	}
}

func (c *Consumer) logWarn(ctx context.Context, message string, attrs ...any) {
	if c.cfg.Logger != nil {
		c.cfg.Logger.WarnContext(ctx, message, attrs...)
	}
}

func (c *Consumer) logError(ctx context.Context, message string, attrs ...any) {
	if c.cfg.Logger != nil {
		c.cfg.Logger.ErrorContext(ctx, message, attrs...)
	}
}

func normalizeConsumerConfig(cfg ConsumerConfig) ConsumerConfig {
	if cfg.MinBytes <= 0 {
		cfg.MinBytes = 10e3
	}
	if cfg.MaxBytes <= 0 {
		cfg.MaxBytes = 10e6
	}
	if cfg.MaxWait <= 0 {
		cfg.MaxWait = 2 * time.Second
	}
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 3
	}
	if cfg.RetryBackoff <= 0 {
		cfg.RetryBackoff = 500 * time.Millisecond
	}
	return cfg
}
