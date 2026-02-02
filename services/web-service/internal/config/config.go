package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Redis  RedisConfig
	Server ServerConfig
	Api    ApiConfig
}

type ServerConfig struct {
	Port string
}

type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

type ApiConfig struct {
	AuthURL string
	DataURL string
	Timeout int
}

func Load() *Config {

	_ = godotenv.Load()

	serverConfig := ServerConfig{
		Port: getEnv("PORT", "8080"),
	}

	redisConfig := RedisConfig{
		URL:      getEnv("REDIS_URL", "redis:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       getIntEnv("REDIS_DB", 0),
	}

	apiConfig := ApiConfig{
		AuthURL: getEnv("AUTH_URL", "localhost:8081"),
		DataURL: getEnv("DATA_URL", "localhost:8082"),
		Timeout: getIntEnv("TIMEOUT", 30),
	}

	return &Config{
		Server: serverConfig,
		Redis:  redisConfig,
		Api:    apiConfig,
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
