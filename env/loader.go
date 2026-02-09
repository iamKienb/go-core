package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using OS env")
	}
}

func Get(key string) string {
	return os.Getenv(key)
}
