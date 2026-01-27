package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	DatabaseURL      string

	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool

	RedisURL      string
	RedisPassword string
	RedisDB       int

	AuthServiceURL string

	Port string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		PostgresHost:     getEnv("POSTGRES_HOST", "postgres"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		PostgresDB:       getEnv("POSTGRES_DB", "postgres"),

		MinioEndpoint:  getEnv("MINIO_ENDPOINT", "minio:9000"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", "minio"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", ""),
		MinioBucket:    getEnv("MINIO_BUCKET", "files"),
		MinioUseSSL:    getBoolValue("MINIO_USE_SSL", false),

		RedisURL: getEnv("REDIS_URL", "redis:6378"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB: getIntEnv("REDIS_DB", 0),

		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://auth-service:8081"),

		Port: getEnv("PORT", "8082"),
	}

	cfg.DatabaseURL = cfg.getDatabaseURL()

	return cfg
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
	value, ok := os.LookupEnv(key)
	if !ok {
		value = defaultValue
	}
	return value
}

func (c *Config) getDatabaseURL() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresDB,
	)
}
