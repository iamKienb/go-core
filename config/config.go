package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type PostgresConfig struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Db              string
	MaxIdleConns    int32
	MaxOpenConns    int32
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database int
	PoolSize int
}

type JwtConfig struct {
	accessSecret  string
	refreshSecret string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

type Server struct {
	Port int
}

type Argon2Config struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  int
	keyLength   uint32
}

type KafkaConfig struct {
	Brokers  []string
	ClientID string
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnvInt(key string) int {
	v := os.Getenv(key)
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("invalid env %s=%s", key, v)
	}
	return i
}

func GetEnvDuration(key string, fallback time.Duration) time.Duration {
	s := GetEnv(key, "")
	if s == "" {
		return fallback
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return fallback
	}
	return d
}
