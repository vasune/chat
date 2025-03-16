package usecases

import (
	"chat/internal/repository"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type ChatUseCase struct {
	clients   map[*websocket.Conn]bool
	broadcast chan Message
	upgrader  websocket.Upgrader
	mu        sync.Mutex
	userRepo  repository.UserRepository
}

func NewChatUseCase(userRepo repository.UserRepository) *ChatUseCase {
	return &ChatUseCase{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan Message),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		userRepo: userRepo,
	}
}

func (uc *ChatUseCase) HandleMessages() {
	for {
		msg := <-uc.broadcast
		uc.mu.Lock()
		for client := range uc.clients {
			if err := client.WriteJSON(msg); err != nil {
				client.Close()
				delete(uc.clients, client)
			}
		}
		uc.mu.Unlock()
	}
}

func (uc *ChatUseCase) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := uc.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	defer ws.Close()

	userID := r.Context().Value("userID").(float64)
	user, err := uc.userRepo.FindByUserID(uint(userID))
	if err != nil {
		log.Printf("Failed to find user by ID: %v", err)
		http.Error(w, "Failed to find user by ID", http.StatusInternalServerError)
		return
	}

	uc.mu.Lock()
	uc.clients[ws] = true
	uc.mu.Unlock()

	for {
		var msg Message
		if err := ws.ReadJSON(&msg); err != nil {
			uc.mu.Lock()
			delete(uc.clients, ws)
			uc.mu.Unlock()
			break
		}
		msg.Username = user.Username
		uc.broadcast <- msg
	}
}
