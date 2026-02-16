package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

type PosetgresConfig struct {
	DB       string
	User     string
	Password string
	Host     string
	Port     string

	DatabaseURL string
}

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type Config struct {
	Postgres PosetgresConfig
	Redis    RedisConfig
	Minio    MinioConfig

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
		Postgres: PosetgresConfig{
			DB:       getEnv("POSTGRES_DB", "db"),
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
		},

		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", "password"),
			DB:       redisDB,
		},
		SessionTTL: time.Duration(sessionTTL) * time.Second,
		Port:       getEnv("PORT", "8081"),
	}
}

func (c *PosetgresConfig) GetDatabaseURL() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DB,
	)
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}
