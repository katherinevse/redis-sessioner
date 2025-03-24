package main

import (
	"app/app/internal/config"
	"app/app/internal/handler"
	cache "app/app/internal/redis"
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Ошибка конфигурации: %v", err)
	}

	//TODO
	//	logger := setupLogger(cfg.LoggerConfig.Level)
	//	logger.Info("Loaded configuration", slog.Any("config", cfg))

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisCfg.Addr,
		Password: cfg.RedisCfg.Password,
		DB:       cfg.RedisCfg.DB,
	})

	if err = redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}
	cacheClient := cache.NewClient(redisClient)
	handler := handler.New(cacheClient) //TODO  service слой

	//postgresDB, err := db.NewPostgresDB(cfg)
	//if err != nil {
	//	log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	//}
	//defer postgresDB.Close()
	//
	//r := mux.NewRouter()
	//
	//handler.RegisterRoutes(r)
	//
	//fmt.Printf("Server is running at %s:%d\n", cfg.ServerCfg.Host, cfg.ServerCfg.Port)
	//fmt.Printf("Redis Address: %s\n", cfg.RedisCfg.Addr)
}
