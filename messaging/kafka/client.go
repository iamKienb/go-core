package kafka

import (
	"sync"

	"github.com/segmentio/kafka-go"
)

type Client struct {
	brokers []string
	writer  map[string]*kafka.Writer
	mu      *sync.Mutex
}

func New(brokers []string) *Client {
	return &Client{
		brokers: brokers,
	}
}

func (c *Client) GetWriter(topic string) *kafka.Writer
