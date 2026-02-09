package redis

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database int
	PoolSize int
}
