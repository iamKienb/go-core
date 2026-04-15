package config

import (
	"time"
)

type PostgresConfig struct {
	Host            string        `env:"USER_COMMAND_SERVICE_DB_HOST"`
	Port            int           `env:"USER_COMMAND_SERVICE_DB_PORT"`
	Username        string        `env:"USER_COMMAND_SERVICE_DB_USERNAME"`
	Password        string        `env:"USER_COMMAND_SERVICE_DB_PASSWORD"`
	Db              string        `env:"USER_COMMAND_SERVICE_DB_NAME"`
	MaxIdleConns    int           `env:"USER_COMMAND_SERVICE_DB_MAX_IDLE_CONNS"`
	MaxOpenConns    int           `env:"USER_COMMAND_SERVICE_DB_MAX_OPEN_CONNS"`
	ConnMaxLifetime time.Duration `env:"USER_COMMAND_SERVICE_DB_CONN_MAX_LIFETIME"`
}

type RedisConfig struct {
	Host     string `env:"USER_COMMAND_SERVICE_REDIS_HOST"`
	Port     int    `env:"USER_COMMAND_SERVICE_REDIS_PORT"`
	Username string `env:"USER_COMMAND_SERVICE_REDIS_USERNAME"`
	Password string `env:"USER_COMMAND_SERVICE_REDIS_PASSWORD"`
	Db       int    `env:"USER_COMMAND_SERVICE_REDIS_DB"`
	PoolSize int    `env:"USER_COMMAND_SERVICE_REDIS_POOL_SIZE"`
}

type JwtConfig struct {
	AccessSecret  string        `env:"USER_COMMAND_SERVICE_ACCESS_SECRET"`
	RefreshSecret string        `env:"USER_COMMAND_SERVICE_REFRESH_SECRET"`
	AccessExpiry  time.Duration `env:"USER_COMMAND_SERVICE_ACCESS_EXPIRY"`
	RefreshExpiry time.Duration `env:"USER_COMMAND_SERVICE_REFRESH_EXPIRY"`
}

type Argon2Config struct {
	Memory      uint32 `env:"USER_COMMAND_SERVICE_ARGON2_MEMORY"`
	Iterations  uint32 `env:"USER_COMMAND_SERVICE_ARGON2_ITERATIONS"`
	Parallelism uint8  `env:"USER_COMMAND_SERVICE_ARGON2_PARALLELISM"`
	SaltLength  int    `env:"USER_COMMAND_SERVICE_ARGON2_SALT_LENGTH"`
	KeyLength   uint32 `env:"USER_COMMAND_SERVICE_ARGON2_KEY_LENGTH"`
}

type KafkaConfig struct {
	Brokers  []string `env:"USER_COMMAND_SERVICE_KAFKA_BROKERS" envSeparator:","`
	ClientID string   `env:"USER_COMMAND_SERVICE_KAFKA_CLIENT_ID"`
}
type Server struct {
	GrpcPort int `env:"USER_COMMAND_SERVICE_GRPC_PORT"`
}
