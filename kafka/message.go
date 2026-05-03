package kafkax

import (
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	HeaderRetryCount      = "x-retry-count"
	HeaderOriginalTopic   = "x-original-topic"
	HeaderOriginalOffset  = "x-original-offset"
	HeaderOriginalGroupID = "x-original-group-id"
	HeaderFailureReason   = "x-failure-reason"
	HeaderFailureAt       = "x-failure-at"
	HeaderIdempotencyKey  = "x-idempotency-key"
)

type Message struct {
	Topic     string
	Key       []byte
	Value     []byte
	Headers   map[string]string
	Partition int
	Offset    int64
	Time      time.Time
}

func (m Message) Header(key string) string {
	if m.Headers == nil {
		return ""
	}
	return m.Headers[strings.ToLower(key)]
}

func (m *Message) SetHeader(key, value string) {
	if m.Headers == nil {
		m.Headers = make(map[string]string)
	}
	m.Headers[strings.ToLower(key)] = value
}

func (m Message) RetryCount() int {
	raw := m.Header(HeaderRetryCount)
	if raw == "" {
		return 0
	}

	count, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}

	return count
}

func (m Message) IdempotencyKey() string {
	if key := strings.TrimSpace(m.Header(HeaderIdempotencyKey)); key != "" {
		return key
	}
	if len(m.Key) > 0 {
		return string(m.Key)
	}
	return ""
}

func toKafkaMessage(msg Message) kafka.Message {
	headers := make([]kafka.Header, 0, len(msg.Headers))
	for key, value := range msg.Headers {
		headers = append(headers, kafka.Header{
			Key:   key,
			Value: []byte(value),
		})
	}

	return kafka.Message{
		Topic:   msg.Topic,
		Key:     msg.Key,
		Value:   msg.Value,
		Headers: headers,
		Time:    msg.Time,
	}
}

func fromKafkaMessage(topic string, msg kafka.Message) Message {
	headers := make(map[string]string, len(msg.Headers))
	for _, header := range msg.Headers {
		headers[strings.ToLower(header.Key)] = string(header.Value)
	}

	return Message{
		Topic:     topic,
		Key:       msg.Key,
		Value:     msg.Value,
		Headers:   headers,
		Partition: msg.Partition,
		Offset:    msg.Offset,
		Time:      msg.Time,
	}
}
