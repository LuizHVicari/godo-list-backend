package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	return Config{
		ServerPort: os.Getenv("SERVER_PORT"),
	}

}
