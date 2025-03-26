package db

import (
	"app/app/internal/config"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Port,
		cfg.DB.Name,
	)
	fmt.Println("DSN:", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	log.Println("Подключено к PostgreSQL")
	return db, nil
}
