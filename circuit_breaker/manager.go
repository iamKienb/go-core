package cbx

import (
	"sync"

	configx "github.com/iamKienb/shopify-go-platform/config"
)

type Manager struct {
	breakers map[string]*CircuitBreaker
	mu       sync.RWMutex
	cfg      configx.CircuitBreakerConfig
}

func NewManager(cfg configx.CircuitBreakerConfig) *Manager {
	return &Manager{
		breakers: make(map[string]*CircuitBreaker),
		cfg:      cfg,
	}
}

func (m *Manager) Get(name string) *CircuitBreaker {
	m.mu.RLock()
	breaker, ok := m.breakers[name]
	m.mu.RUnlock()
	if ok {
		return breaker
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if breaker, ok = m.breakers[name]; ok {
		return breaker
	}

	newBreaker := NewCircuitBreaker(name, m.cfg)
	m.breakers[name] = newBreaker

	return newBreaker
}
