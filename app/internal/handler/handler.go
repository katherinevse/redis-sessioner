package handler

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type SessionManager interface {
	GetSession(ctx context.Context, sessionID string) (string, error)
}

type Handler struct {
	sessionService SessionManager
}

func New(sessionService SessionManager) *Handler {
	return &Handler{
		sessionService: sessionService,
	}
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	sessionID := mux.Vars(r)["sessionID"] //мапа приходит проверить все данные

	sessionData, err := h.sessionService.GetSession(r.Context(), sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if sessionData == "" {
		http.Error(w, "Сессия не найдена", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"sessionData": sessionData})
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/session/{sessionID}", h.GetSession).Methods("GET")
	//r.HandleFunc("/session/{sessionID}", h.DeleteSession).Methods("DELETE")
	//r.HandleFunc("/session", h.UpdateSession).Methods("PUT")

}
