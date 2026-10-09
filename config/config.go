package configx

import (
	"time"

	"github.com/segmentio/kafka-go"
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

type ElasticSearchConfig struct {
	Addresses []string `env:"_ELASTICSEARCH_ADDRESSES"`
	Username  string   `env:"_ELASTICSEARCH_USERNAME"`
	Password  string   `env:"_ELASTICSEARCH_PASSWORD"`
	CloudID   string   `env:"_ELASTICSEARCH_CLOUD_ID"`
	APIKey    string   `env:"_ELASTICSEARCH_API_KEY"`
}

type JwtConfig struct {
	AccessExpiry   time.Duration `env:"_ACCESS_EXPIRY"`
	RefreshExpiry  time.Duration `env:"_REFRESH_EXPIRY"`
	PrivateKeyPath string        `env:"_PRIVATE_KEY_PATH" envDefault:"private.pem"`
	PublicKeyPath  string        `env:"_PUBLIC_KEY_PATH" envDefault:"public.pem"`
}

type Argon2Config struct {
	Memory      uint32 `env:"_ARGON2_MEMORY"`
	Iterations  uint32 `env:"_ARGON2_ITERATIONS"`
	Parallelism uint8  `env:"_ARGON2_PARALLELISM"`
	SaltLength  int    `env:"_ARGON2_SALT_LENGTH"`
	KeyLength   uint32 `env:"_ARGON2_KEY_LENGTH"`
}

type KafkaConfig struct {
	Brokers      []string      `env:"_KAFKA_BROKERS"`
	ClientID     string        `env:"_KAFKA_CLIENT_ID"`
	DialTimeout  time.Duration `env:"_KAFKA_DIAL_TIMEOUT" envDefault:"5s"`
	ReadTimeout  time.Duration `env:"_KAFKA_READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout time.Duration `env:"_KAFKA_WRITE_TIMEOUT" envDefault:"10s"`
}

type ProducerConfig struct {
	Topic          string             `env:"_PRODUCER_TOPIC"`
	Balancer       string             `env:"_PRODUCER_BALANCER"`
	Compression    string             `env:"_PRODUCER_COMPRESSION"`
	BatchTimeout   time.Duration      `env:"_PRODUCER_BATCH_TIMEOUT"`
	BatchBytes     int64              `env:"_PRODUCER_BATCH_BYTES"`
	RequiredAcks   kafka.RequiredAcks `env:"_PRODUCER_REQUIRED_ACKS"`
	AllowAutoTopic bool               `env:"_PRODUCER_ALLOW_AUTO_TOPIC"`
	WriteTimeout   time.Duration      `env:"_PRODUCER_WRITE_TIMEOUT"`
	ReadTimeout    time.Duration      `env:"_PRODUCER_READ_TIMEOUT"`
	MaxAttempts    int                `env:"_PRODUCER_MAX_ATTEMPTS"`
}

type ConsumerConfig struct {
	GroupID      string `env:"_CONSUMER_GROUP_ID"`
	Topics       []string
	DLQTopic     string        `env:"_CONSUMER_DLQ_TOPIC"`
	MinBytes     int           `env:"_CONSUMER_MIN_BYTES"`
	MaxBytes     int           `env:"_CONSUMER_MAX_BYTES"`
	MaxWait      time.Duration `env:"_CONSUMER_MAX_WAIT"`
	MaxAttempts  int           `env:"_CONSUMER_MAX_ATTEMPTS"`
	RetryBackoff time.Duration `env:"_CONSUMER_RETRY_BACKOFF"`
}

type CircuitBreakerConfig struct {
	MaxRequests  uint32        `env:"_BREAKER_MAX_REQUESTS"`
	Interval     time.Duration `env:"_BREAKER_INTERVAL"`
	Timeout      time.Duration `env:"_BREAKER_TIMEOUT"`
	FailureRatio float64       `env:"_BREAKER_FAILURE_RATIO"`
	ThresholdCnt uint32        `env:"_BREAKER_THRESHOLD_CNT"`
}

type TelemetryConfig struct {
	JaegerEndpoint string `env:"_JAEGER_ENDPOINT"`
	MetricsPath    string `env:"_METRICS_PATH"`
	Enabled        bool   `env:"_TELEMETRY_ENABLED"`
}

type Server struct {
	GrpcPort          int    `env:"_GRPC_PORT"`
	ApiGatewayAddress string `env:"_API_GATEWAY_ADDRESS"`
}
