package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using OS env")
	}
}

func GetString(key string) string {
	return os.Getenv(key)
}

func GetInt(key string) int {
	v := os.Getenv(key)
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("invalid env %s=%s", key, v)
	}
	return i
}
