package configx

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func Loader[T any]() (*T, error) {
	_ = godotenv.Load()
	var cfg T
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("config: failed to parse struct %T: %w", cfg, err)
	}
	return &cfg, nil
}
