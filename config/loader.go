package configx

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func Loader[T any]() (*T, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("config: .env file not found, using system environment")
	}
	var cfg T
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("config: failed to parse struct %T: %w", cfg, err)
	}
	return &cfg, nil
}
