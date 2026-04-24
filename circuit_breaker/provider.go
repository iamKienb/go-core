package cb

import (
	"errors"
	"fmt"

	configx "github.com/iamKienb/shopify-go-platform/config"
	"github.com/sony/gobreaker"
)

type CircuitBreaker struct {
	cb *gobreaker.CircuitBreaker
}

func NewCircuitBreaker(cfg configx.CircuitBreakerConfig) *CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        cfg.Name,
		MaxRequests: cfg.MaxRequests,
		Interval:    cfg.Interval,
		Timeout:     cfg.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			ratio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= cfg.ThresholdCnt && ratio >= cfg.FailureRatio
		},
	}

	return &CircuitBreaker{
		cb: gobreaker.NewCircuitBreaker(settings),
	}
}

func (c *CircuitBreaker) Execute(fn func() (interface{}, error)) (interface{}, error) {
	result, err := c.cb.Execute(fn)
	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			return nil, fmt.Errorf("circuit breaker is open: service unavailable")
		}
		return nil, err
	}

	return result, nil
}
