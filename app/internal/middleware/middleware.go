package middleware

import (
	"context"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Incr(ctx context.Context, key string) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
}

type RateLimiter struct {
	client   RedisClient
	rate     int
	interval time.Duration
}

func NewRateLimiter(client RedisClient, rate int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		client:   client,
		rate:     rate,
		interval: interval,
	}
}

func (rl *RateLimiter) AllowRequest(ctx context.Context, key string) bool {
	count, err := rl.client.Incr(ctx, key).Result()
	if err != nil {
		return false
	}

	if count == 1 {
		rl.client.Expire(ctx, key, rl.interval)
	}

	return count <= int64(rl.rate)
}

func RateLimitMiddleware(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ip := r.RemoteAddr // Можно использовать API ключи или userID

			if !rl.AllowRequest(ctx, ip) {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
