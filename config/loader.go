package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using OS env")
	}

	return &Config{
		Postgres: PostgresConfig{
			Host:            getEnv("PG_HOST"),
			Port:            getEnvInt("PG_PORT"),
			Username:        getEnv("PG_USER"),
			Password:        getEnv("PG_PASS"),
			Db:              getEnv("PG_DB"),
			MaxIdleConns:    int32(getEnvInt("PG_MAX_IDLE_CONNS")),
			MaxOpenConns:    int32(getEnvInt("PG_MAX_OPEN_CONNS")),
			ConnMaxLifetime: time.Duration(getEnvInt("PG_CONN_MAX_LIFE_TIME")) * time.Second,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST"),
			Port:     getEnvInt("REDIS_PORT"),
			Username: getEnv("REDIS_USERNAME"),
			Password: getEnv("REDIS_PASSWORD"),
			Database: getEnvInt("REDIS_DB"),
			PoolSize: getEnvInt("REDIS_POOL_SIZE"),
		},

		Kafka: KafkaConfig{
			Brokers: []string{
				getEnv("KAFKA_BROKER_1"),
				getEnv("KAFKA_BROKER_2"),
				getEnv("KAFKA_BROKER_3"),
			},
			ClientID: "go-platform",
		},
	}
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func getEnvInt(key string) int {
	v := os.Getenv(key)
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("invalid env %s=%s", key, v)
	}
	return i
}
