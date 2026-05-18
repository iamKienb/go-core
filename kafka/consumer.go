package kafkax

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	configx "github.com/iamKienb/go-core/config"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader      *kafka.Reader
	handler     ConsumerHandler
	dlqProducer *Producer
	cfg         configx.ConsumerConfig
	logger      *slog.Logger
}

func NewConsumer(service KafkaXService, cfg configx.ConsumerConfig, logger *slog.Logger, handler ConsumerHandler) (*Consumer, error) {
	if service == nil || handler == nil {
		return nil, errors.New("kafka consumer: client service and handler must not be nil")
	}

	if len(cfg.Topics) == 0 {
		return nil, errors.New("kafka consumer: topics list must not be empty")
	}
	if strings.TrimSpace(cfg.GroupID) == "" {
		return nil, errors.New("kafka consumer: group id must not be empty")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     service.Brokers(),
		GroupID:     cfg.GroupID,
		GroupTopics: cfg.Topics,
		MinBytes:    cfg.MinBytes,
		MaxBytes:    cfg.MaxBytes,
		MaxWait:     cfg.MaxWait,
		Dialer:      service.Dialer(),
	})

	var dlqProducer *Producer
	if cfg.DLQTopic != "" {
		var err error
		dlqProducer, err = NewProducer(service, configx.ProducerConfig{
			Topic:          cfg.DLQTopic,
			Balancer:       "least_bytes",
			BatchTimeout:   50 * time.Millisecond,
			AllowAutoTopic: true,
			MaxAttempts:    3,
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
		logger:      logger,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	c.logInfo(ctx, "kafka consumer starting")
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

	msg := fromKafkaMessage(kmsg.Topic, kmsg)

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
				slog.String("topic", msg.Topic),
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
	dlqMsg.SetHeader(HeaderOriginalTopic, msg.Topic)
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
	if c.logger != nil {
		c.logger.InfoContext(ctx, msg, args...)
	}
}

func (c *Consumer) logWarn(ctx context.Context, message string, attrs ...any) {
	if c.logger != nil {
		c.logger.WarnContext(ctx, message, attrs...)
	}
}

func (c *Consumer) logError(ctx context.Context, message string, attrs ...any) {
	if c.logger != nil {
		c.logger.ErrorContext(ctx, message, attrs...)
	}
}
