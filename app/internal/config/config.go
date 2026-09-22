package config

import "os"

type Config struct {
	Addr     string
	LogLevel string
}

func Load() (*Config, error) {
	c := &Config{
		Addr: getenv("APP_ADDR", ":8080"),
		LogLevel: getenv("APP_LOG_LEVEL", "info"),
	}
	return c, nil
}

func getenv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
