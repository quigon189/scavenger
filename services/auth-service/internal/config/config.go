package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	RedisURL    string
	RedisTTL    string
	JWTSecret   string
	Port        int
}

func Load() *Config {
	_ = godotenv.Load()
	port, err := strconv.Atoi(getEnv("PORT", "8081"))
	if err != nil {
		port = 8081
	}

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgresql://scavenger:scavenger123@localhost:5432/scavenger?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "localhost:6379"),
		RedisTTL:    getEnv("REDIS_TTL", "24h"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-jwt"),
		Port:        port,
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}
