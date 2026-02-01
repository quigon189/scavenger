package config

import (
	"os"
	"strconv"
)

type Config struct {
	Redis  RedisConfig
	Server ServerConfig
}

type ServerConfig struct {
	Port string
}

type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

func Load() *Config {
	serverConfig := ServerConfig{
		Port: getEnv("PORT", "8080"),
	}

	redisConfig := RedisConfig{
		URL:      getEnv("REDIS_URL", "redis:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       getIntEnv("REDIS_DB", 0),
	}

	return &Config{
		Server: serverConfig,
		Redis:  redisConfig,
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}

	return value
}

func getIntEnv(key string, defaultValue int) int {
	strValue, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	value, err := strconv.Atoi(strValue)
	if err != nil {
		return defaultValue
	}

	return value
}
