package service

type Session struct {
	dbRepo      SessionRepository
	redisClient RedisClient
}
