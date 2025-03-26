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
	DB        Database
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

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

func New() (*Config, error) {
	err := godotenv.Load("internal/config/.env")
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

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("не удалось преобразовать DB_PORT: %w", err)
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
		DB: Database{
			Host:     os.Getenv("DB_HOST"),
			Port:     dbPort,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		LogCfg: LogConfig{
			Level: os.Getenv("LOG_LEVEL"),
		},
	}, nil
}
