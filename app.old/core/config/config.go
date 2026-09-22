package config

import (
	"fmt"
	"log"
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

	SignedURLTTL int
}

type Config struct {
	Postgres PosetgresConfig
	Redis    RedisConfig
	Minio    MinioConfig

	CookieSecret string
	SessionTTL   time.Duration
	Port         string
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
		Minio: MinioConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "minio:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minio"),
			SecretKey: getEnv("MINIO_SECRET_KEY", ""),
			Bucket:    getEnv("MINIO_BUCKET", "files"),
			UseSSL:    getBoolValue("MINIO_USE_SSL", false),
			SignedURLTTL: getIntEnv("MINIO_SIGNED_URL_TTL", 300),
		},
		CookieSecret:  getEnv("COOKIE_SECRET", "secret-key"),
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

func getIntEnv(key string, defaultValue int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("WARN Failed to parse int env %s: %v", key, err)
		return defaultValue
	}

	return intValue
}

func getBoolValue(key string, defaultValue bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("WARN Failed to parse bool env %s: %v", key, err)
		return defaultValue
	}

	return boolValue
}


func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}
