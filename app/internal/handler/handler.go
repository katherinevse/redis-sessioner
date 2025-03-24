package handler

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Handler struct {
	redisClient RedisClient
}

func New(redisClient RedisClient) *Handler {
	return &Handler{
		redisClient: redisClient,
	}
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID") //извлекаем из запроса
	ttl := time.Duration(30) * time.Minute
	key := fmt.Sprintf("session:%s", userID)
	value := "active"

	err := h.redisClient.Set(context.Background(), key, value, ttl)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка создания сессии: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Сессия создана"))
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/session/create", h.CreateSession).Methods("POST")
	router.HandleFunc("/session/{userID}", h.GetSession).Methods("GET")
	router.HandleFunc("/session/delete", h.DeleteSession).Methods("DELETE")
}
