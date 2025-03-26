package repository

import (
	"context"
	"database/sql"
)

type PostgresRepository interface {
	//ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	//QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type PostgresRepo struct {
	db PostgresRepository
}

func New(db PostgresRepository) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (r *PostgresRepo) GetSession(ctx context.Context, sessionID string) (string, error) {
	var data string
	err := r.db.QueryRowContext(ctx, "SELECT data FROM sessions WHERE id = $1", sessionID).Scan(&data)
	if err != nil {
		return "", err
	}
	return data, nil

}
