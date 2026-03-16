package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseHost     string
	DatabasePort     int
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabaseSSLMode  string
	CacheHost        string
	CachePort        int
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	return Config{
		DatabaseHost:     os.Getenv("DB_HOST"),
		DatabasePort:     mustAtoi("DB_PORT"),
		DatabaseUser:     os.Getenv("DB_USER"),
		DatabasePassword: os.Getenv("DB_PASSWORD"),
		DatabaseName:     os.Getenv("DB_NAME"),
		DatabaseSSLMode:  os.Getenv("DB_SSLMODE"),
		CacheHost:        os.Getenv("CACHE_HOST"),
		CachePort:        mustAtoi("CACHE_PORT"),
	}

}

func mustAtoi(key string) int {
	value := os.Getenv(key)
	parsed, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("invalid %s: %q", key, value))
	}

	return parsed
}
