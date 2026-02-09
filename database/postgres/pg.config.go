package postgres

import "time"

type Config struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Db              string
	MaxIdleConns    int32
	MaxOpenConns    int32
	ConnMaxLifetime time.Duration
}
