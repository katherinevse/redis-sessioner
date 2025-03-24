package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServerCfg ServerConfig
	RedisCfg  RedisConfig
	LogCfg    LogConfig
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type RedisConfig struct {
	Addr       string `yaml:"addr"`
	Password   string `yaml:"password"`
	DB         int    `yaml:"db"`
	SessionTTL int64  `yaml:"session_ttl"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

func New() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	} else {
		fmt.Println("Файл .env успешно загружен")
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, fmt.Errorf("не удалось преобразовать REDIS_DB: %w", err)
	}

	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		return nil, fmt.Errorf("не удалось преобразовать SERVER_PORT: %w", err)
	}

	return &Config{
		ServerCfg: ServerConfig{
			Host: os.Getenv("SERVER_HOST"),
			Port: serverPort,
		},
		RedisCfg: RedisConfig{
			Addr:       os.Getenv("REDIS_ADDR"),
			Password:   os.Getenv("REDIS_PASS"),
			DB:         redisDB,
			SessionTTL: 1800,
		},
		LogCfg: LogConfig{
			Level: os.Getenv("LOG_LEVEL"),
		},
	}, nil
}
