package cache

import (
	"github.com/redis/go-redis/v9"

	"github.com/arrebole/vibe-coding-starter/internal/config"
)

func NewRedisClient(cfg config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}
