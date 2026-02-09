package redis

type Config struct {
	Address  string
	Username string
	Password string
	Database int
	PoolSize int
}
