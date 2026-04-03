package kafka

import "time"

type Message struct {
	Key       []byte
	Value     []byte
	Headers   map[string]string
	Timestamp time.Time
}
