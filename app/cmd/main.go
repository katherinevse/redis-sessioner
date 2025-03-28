package main

import (
	"app/internal/config"
	"app/internal/handler"
	"app/internal/middleware"
	cache "app/internal/redis"
	"app/internal/repository"
	"app/internal/service"
	"app/pkg/db"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Ошибка конфигурации: %v", err)
	}

	logger := setupLogger(cfg.LogCfg.Level)
	logger.Info("Loaded configuration", slog.Any("config", cfg))

	dbConn, err := db.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	defer dbConn.Close()

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

	rateLimiter := middleware.NewRateLimiter(rawRedisClient, 5, time.Minute)

	r := mux.NewRouter()
	r.Use(middleware.RateLimitMiddleware(rateLimiter))

	h.RegisterRoutes(r)

	serverAddr := fmt.Sprintf("%s:%d", cfg.ServerCfg.Host, cfg.ServerCfg.Port)
	fmt.Printf("Server is running at %s\n", serverAddr)
	fmt.Printf("Redis Address: %s\n", cfg.RedisCfg.Addr)

	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

}

// TODO добавить логи везде!
func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}
