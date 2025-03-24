package main

import (
	"app/app/internal/config"
	cache "app/app/internal/redis"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Ошибка конфигурации: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisCfg.Addr,
		Password: cfg.RedisCfg.Password,
		DB:       cfg.RedisCfg.DB,
	})

	if err = redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}
	cacheClient := cache.NewClient(redisClient)

	fmt.Printf("Server is running at %s:%d\n", cfg.ServerCfg.Host, cfg.ServerCfg.Port)
	fmt.Printf("Redis Address: %s\n", cfg.RedisCfg.Addr)
}
