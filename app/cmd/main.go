package main

import (
	"app/app/internal/config"
	"app/app/internal/handler"
	cache "app/app/internal/redis"
	"app/app/internal/repository"
	"app/app/internal/service"
	"app/app/pkg/db"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Ошибка конфигурации: %v", err)
	}

	dbConn, err := db.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	defer dbConn.Close()

	//newsql := sql.DB{}
	//newsql.QueryContext()

	//TODO
	//	logger := setupLogger(cfg.LoggerConfig.Level)
	//	logger.Info("Loaded configuration", slog.Any("config", cfg))

	rawRedisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisCfg.Addr,
		Password: cfg.RedisCfg.Password,
		DB:       cfg.RedisCfg.DB,
	})

	if err = rawRedisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}

	redisClient := cache.NewClient(rawRedisClient)
	postgresRepo := repository.New(dbConn)
	sessionService := service.New(postgresRepo, redisClient)
	h := handler.New(sessionService) //TODO  service слой

	r := mux.NewRouter()
	h.RegisterRoutes(r)

	fmt.Printf("Server is running at %s:%d\n", cfg.ServerCfg.Host, cfg.ServerCfg.Port)
	fmt.Printf("Redis Address: %s\n", cfg.RedisCfg.Addr)
}
