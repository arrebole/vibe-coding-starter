package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv           string
	HTTPAddr         string
	LogLevel         string
	DatabaseDSN      string
	Redis            RedisConfig
	ExternalDemoGRPC GRPCClientConfig
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type GRPCClientConfig struct {
	Addr    string
	Timeout time.Duration
}

func MustLoad() Config {
	_ = godotenv.Load()

	return Config{
		AppEnv:      getEnv("APP_ENV", "local"),
		HTTPAddr:    getEnv("HTTP_ADDR", ":8080"),
		LogLevel:    getEnv("LOG_LEVEL", "debug"),
		DatabaseDSN: mustGetEnv("DATABASE_DSN"),
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		ExternalDemoGRPC: GRPCClientConfig{
			Addr:    getEnv("EXTERNAL_DEMO_GRPC_ADDR", ""),
			Timeout: time.Duration(getEnvInt("EXTERNAL_DEMO_GRPC_TIMEOUT_MS", 1000)) * time.Millisecond,
		},
	}
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("缺少必要环境变量: %s", key)
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("环境变量 %s 必须是整数: %v", key, err)
	}
	return parsed
}
