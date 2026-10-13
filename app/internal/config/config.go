package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr     string
	LogLevel string
	DBDSN    string

	StorageDriver string // fs or s3
	StorageFSRoot string
	MaxUploadSize int64
}

func Load() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	c := &Config{
		Addr:     getenv("APP_ADDR", ":8080"),
		LogLevel: getenv("APP_LOG_LEVEL", "info"),
		DBDSN:    getenv("DB_DSN", "./data/app.db"),

		StorageDriver: getenv("STORAGE_DRIVER", "fs"),
		StorageFSRoot: getenv("STORAGE_FS_ROOT", "./data"),
		MaxUploadSize: getenvInt64("STORAGE_MAX_UPLOAD_BYTES", 10<<20),
	}
	return c, nil
}

func getenv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func getenvInt64(key string, def int64) int64 {
	if strVal := os.Getenv(key); strVal != "" {
		val, err := strconv.ParseInt(strVal, 10, 64)
		if err == nil {
			return val
		}
	}
	return def
}
