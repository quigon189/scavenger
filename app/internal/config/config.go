package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr     string
	LogLevel string
	DBDSN      string
}

func Load() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	c := &Config{
		Addr:     getenv("APP_ADDR", ":8080"),
		LogLevel: getenv("APP_LOG_LEVEL", "info"),
		DBDSN: getenv("DB_DSN", "./data/app.db"),
	}
	return c, nil
}

func getenv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
