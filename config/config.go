package configx

import (
	"time"
)

type PostgresConfig struct {
	Host            string        `env:"_PG_HOST"`
	Port            int           `env:"_PG_PORT"`
	Username        string        `env:"_PG_USERNAME"`
	Password        string        `env:"_PG_PASSWORD"`
	Db              string        `env:"_PG_NAME"`
	MaxIdleConns    int           `env:"_PG_MAX_IDLE_CONNS"`
	MaxOpenConns    int           `env:"_PG_MAX_OPEN_CONNS"`
	ConnMaxLifetime time.Duration `env:"_PG_CONN_MAX_LIFETIME"`
}

type RedisConfig struct {
	Host     string `env:"_REDIS_HOST"`
	Port     int    `env:"_REDIS_PORT"`
	Username string `env:"_REDIS_USERNAME"`
	Password string `env:"_REDIS_PASSWORD"`
	Db       int    `env:"_REDIS_DB"`
	PoolSize int    `env:"_REDIS_POOL_SIZE"`
}

type JwtConfig struct {
	AccessExpiry  time.Duration `env:"_ACCESS_EXPIRY"`
	RefreshExpiry time.Duration `env:"_REFRESH_EXPIRY"`
}

type Argon2Config struct {
	Memory      uint32 `env:"_ARGON2_MEMORY"`
	Iterations  uint32 `env:"_ARGON2_ITERATIONS"`
	Parallelism uint8  `env:"_ARGON2_PARALLELISM"`
	SaltLength  int    `env:"_ARGON2_SALT_LENGTH"`
	KeyLength   uint32 `env:"_ARGON2_KEY_LENGTH"`
}

type KafkaConfig struct {
	Brokers  []string `env:"_KAFKA_BROKERS" envSeparator:","`
	ClientID string   `env:"_KAFKA_CLIENT_ID"`
}

type CircuitBreakerConfig struct {
	MaxRequests  uint32        `env:"_BREAKER_MAX_REQUESTS"`
	Interval     time.Duration `env:"_BREAKER_INTERVAL"`
	Timeout      time.Duration `env:"_BREAKER_TIMEOUT"`
	FailureRatio float64       `env:"_BREAKER_FAILURE_RATIO"`
	ThresholdCnt uint32        `env:"_BREAKER_THRESHOLD_CNT"`
}

type Server struct {
	GrpcPort int `env:"_GRPC_PORT"`
}
