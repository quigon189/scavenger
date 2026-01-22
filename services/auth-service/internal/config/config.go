package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	PosetgresDB       string
	PostgresUser      string
	PosetgresPassword string
	PostgresHost      string
	PosetgresPort     string

	DatabaseURL string

	RedisURL      string
	RedisPassword string
	RedisDB       int

	SessionTTL time.Duration
	Port       string
}

func Load() *Config {
	_ = godotenv.Load()
	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		redisDB = 0
	}

	sessionTTL, err := strconv.Atoi(getEnv("SESSION_TTL", "86400"))
	if err != nil {
		sessionTTL = 86400
	}

	return &Config{
		DatabaseURL:       getEnv("DATABASE_URL", "postgresql://scavenger:scavenger123@localhost:5432/scavenger?sslmode=disable"),
		PosetgresDB:       getEnv("POSTGRES_DB", "db"),
		PostgresHost:      getEnv("POSTGRES_HOST", "localhost"),
		PosetgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:      getEnv("POSTGRES_USER", "postgres"),
		PosetgresPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		RedisURL:          getEnv("REDIS_URL", "localhost:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", "password"),
		RedisDB:           redisDB,
		SessionTTL:        time.Duration(sessionTTL) * time.Second,
		Port:              getEnv("PORT", "8081"),
	}
}

func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.PostgresUser,
		c.PosetgresPassword,
		c.PostgresHost,
		c.PosetgresPort,
		c.PosetgresDB,
	)
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}
