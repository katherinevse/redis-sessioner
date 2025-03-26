package handler

import (
	"context"
)

type SessionManager interface {
	GetSession(ctx context.Context, sessionID string) (string, error)
}
