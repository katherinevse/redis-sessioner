package service

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

//type Sessioner interface {
//	GetSession(ctx context.Context, sessionID string) (string, error)     // Получить сессию (из Redis или БД)
//	SaveSession(ctx context.Context, sessionID string, data string) error // Сохранить сессию (в БД и Redis)
//	DeleteSession(ctx context.Context, sessionID string) error            // Удалить сессию (из Redis и БД)
//}

type SessionService struct {
	db    DBRepository
	cache RedisClient
}

func New(db DBRepository, cache RedisClient) *SessionService {
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
	data := ""
	err = s.db.QueryRowContext(ctx, "SELECT data FROM sessions WHERE id = $1", sessionID).Scan(&data)
	if err != nil {
		return "", err
	}

	// Кешируем в Redis на 30 минут
	s.cache.Set(ctx, sessionID, data, 30*time.Minute)

	return data, nil
}

// SetSession Сохранить сессию (в БД и Redis)
func (s *SessionService) SetSession(ctx context.Context, sessionID string, data string) error {

}

//Удалить сессию (из Redis и БД)
