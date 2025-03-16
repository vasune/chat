package handlers

import (
	"chat/internal/usecases"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	AuthUC usecases.Auth
}

func NewAuthHandler(useCases usecases.Auth) *AuthHandler {
	return &AuthHandler{AuthUC: useCases}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.AuthUC.SignUp(request.Username, request.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return
	}

	token, err := h.AuthUC.SignIn(request.Username, request.Password)
	if err != nil {
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

type ChatHandler struct {
	ChatUC usecases.Chat
}

func NewChatHandler(useCases usecases.Chat) *ChatHandler {
	return &ChatHandler{ChatUC: useCases}
}

func (h *ChatHandler) HandleMessages() {
	h.ChatUC.HandleMessages()
}

func (h *ChatHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	h.ChatUC.HandleConnections(w, r)
}
