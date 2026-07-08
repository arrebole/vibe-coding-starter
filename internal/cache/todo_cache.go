package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/dto"
)

const todoListCacheKey = "todos:list"

type TodoCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewTodoCache(client *redis.Client, ttl time.Duration) *TodoCache {
	return &TodoCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *TodoCache) GetList(ctx context.Context) ([]dto.TodoResponse, bool, error) {
	payload, err := c.client.Get(ctx, todoListCacheKey).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	var todos []dto.TodoResponse
	if err := json.Unmarshal(payload, &todos); err != nil {
		return nil, false, err
	}
	return todos, true, nil
}

func (c *TodoCache) SetList(ctx context.Context, todos []dto.TodoResponse) error {
	payload, err := json.Marshal(todos)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, todoListCacheKey, payload, c.ttl).Err()
}

func (c *TodoCache) DeleteList(ctx context.Context) error {
	return c.client.Del(ctx, todoListCacheKey).Err()
}
