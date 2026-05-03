package kafkax

import (
	"context"
	"errors"
	"fmt"
	"strings"

	configx "github.com/iamKienb/shopify-go-platform/config"
	"github.com/segmentio/kafka-go"
)

func (x *KafkaX) Ping(ctx context.Context) error {
	conn, err := x.dialer.DialContext(ctx, "tcp", x.cfg.Brokers[0])
	if err != nil {
		return fmt.Errorf("kafka dial failed: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Brokers(); err != nil {
		return fmt.Errorf("kafka metadata fetch failed: %w", err)
	}

	return nil
}

func (x *KafkaX) Dialer() *kafka.Dialer {
	return x.dialer
}

func (x *KafkaX) Transport() *kafka.Transport {
	return x.transport
}

func (x *KafkaX) Brokers() []string {
	return append([]string(nil), x.cfg.Brokers...)
}

func (x *KafkaX) ClientID() string {
	return x.cfg.ClientID
}

func (x *KafkaX) Close() error {
	if x.transport != nil {
		x.transport.CloseIdleConnections()
	}
	return nil
}

func validateConfig(cfg configx.KafkaConfig) error {
	if len(cfg.Brokers) == 0 {
		return errors.New("kafka config: brokers must not be empty")
	}

	for _, broker := range cfg.Brokers {
		if strings.TrimSpace(broker) == "" {
			return errors.New("kafka config: broker contains empty value")
		}
	}

	if strings.TrimSpace(cfg.ClientID) == "" {
		return errors.New("kafka config: client id must not be empty")
	}

	if cfg.DialTimeout <= 0 {
		return errors.New("kafka config: dial timeout must be greater than zero")
	}
	if cfg.ReadTimeout <= 0 {
		return errors.New("kafka config: read timeout must be greater than zero")
	}
	if cfg.WriteTimeout <= 0 {
		return errors.New("kafka config: write timeout must be greater than zero")
	}

	return nil
}
