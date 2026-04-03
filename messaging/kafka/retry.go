package kafka

import (
	"context"
	"fmt"
	"time"
)

type RetryHandler struct {
	maxRetry  int
	baseDelay time.Duration
	publisher *Publisher
}

func NewRetryHandler(max int, pub *Publisher) *RetryHandler {
	return &RetryHandler{
		maxRetry:  max,
		publisher: pub,
	}
}

func (r *RetryHandler) GetRetryCount(msg Message) int {
	retry := 0

	if val, ok := msg.Headers["x-retry"]; ok {
		fmt.Sscanf(val, "%d", &retry)
	}

	return retry
}

func (r *RetryHandler) CanRetry(msg Message) bool {
	return r.GetRetryCount(msg) < r.maxRetry
}

func (r *RetryHandler) Handle(ctx context.Context, msg Message) error {
	retry := r.GetRetryCount(msg) + 1

	delay := r.baseDelay * time.Duration(1<<retry)

	if msg.Headers == nil {
		msg.Headers = make(map[string]string)
	}

	msg.Headers["x-retry"] = fmt.Sprintf("%d", retry)
	msg.Headers["x-delay"] = delay.String()

	return r.publisher.Publish(ctx, msg)
}
