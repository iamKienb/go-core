package config

import "time"

type Config struct {
	Postgres PostgresConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
}

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

type KafkaConfig struct {
	Brokers  []string
	ClientID string
}
