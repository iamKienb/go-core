package cbx

import (
	"errors"

	app_error "github.com/iamKienb/go-core/app_error"
	configx "github.com/iamKienb/go-core/config"
	"github.com/sony/gobreaker"
)

type CircuitBreaker struct {
	cb *gobreaker.CircuitBreaker
}

func NewCircuitBreaker(name string, cfg configx.CircuitBreakerConfig) *CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        name,
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
			return nil, app_error.Unavailable("service temporarily unavailable", err)
		}
		return nil, app_error.Unavailable("service temporarily unavailable", err)
	}

	return result, nil
}
