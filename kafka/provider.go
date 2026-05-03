package kafkax

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	configx "github.com/iamKienb/shopify-go-platform/config"
	"github.com/segmentio/kafka-go"
)

type Client struct {
	cfg       configx.KafkaConfig
	dialer    *kafka.Dialer
	transport *kafka.Transport
}

func New(cfg configx.KafkaConfig) (*Client, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	dialer := &kafka.Dialer{
		ClientID:  cfg.ClientID,
		Timeout:   cfg.DialTimeout,
		DualStack: true,
	}

	transport := &kafka.Transport{
		ClientID: cfg.ClientID,
		Dial: (&net.Dialer{
			Timeout: cfg.DialTimeout,
		}).DialContext,
	}

	client := &Client{
		cfg:       cfg,
		dialer:    dialer,
		transport: transport,
	}

	if err := client.Ping(context.Background()); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Ping(ctx context.Context) error {
	conn, err := c.dialer.DialContext(ctx, "tcp", c.cfg.Brokers[0])
	if err != nil {
		return fmt.Errorf("kafka dial failed: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Brokers(); err != nil {
		return fmt.Errorf("kafka metadata fetch failed: %w", err)
	}

	return nil
}

func (c *Client) Dialer() *kafka.Dialer {
	return c.dialer
}

func (c *Client) Transport() *kafka.Transport {
	return c.transport
}

func (c *Client) Brokers() []string {
	return append([]string(nil), c.cfg.Brokers...)
}

func (c *Client) ClientID() string {
	return c.cfg.ClientID
}

func (c *Client) Close() error {
	if c.transport != nil {
		c.transport.CloseIdleConnections()
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
