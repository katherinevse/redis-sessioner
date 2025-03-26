package service

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type DBRepo interface {
	GetSession(ctx context.Context, sessionID string) (string, error)
}

type SessionService struct {
	db    DBRepo
	cache RedisClient
}

func New(db DBRepo, cache RedisClient) *SessionService {
	return &SessionService{db: db, cache: cache}
}

// GetSession Получить сессию (из Redis или БД)
func (s *SessionService) GetSession(ctx context.Context, sessionID string) (string, error) {
	cachedValue, err := s.cache.Get(ctx, sessionID).Result()
	if err == nil {
		return cachedValue, nil
	} else if !errors.Is(redis.Nil, err) {
		return "", err
	}

	// Если нет в Redis — берём из БД
	data, err := s.db.GetSession(ctx, sessionID)
	if err != nil {
		return "", err
	}

	// Кешируем в Redis на 30 минут
	s.cache.Set(ctx, sessionID, data, 30*time.Minute)

	return data, nil
}
