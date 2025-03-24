package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type Client struct {
	client RedisClient
}

func NewClient(client RedisClient) *Client {
	return &Client{
		client: client,
	}
}

// Set сохраняет значение в Redis с указанным TTL
func (r *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(ctx, key, value, expiration)
}

// Get получает значение из Redis
func (r *Client) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Get(ctx, key)
}

// Delete удаляет ключ из Redis
func (r *Client) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.client.Del(ctx, keys...)
}
