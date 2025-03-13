package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"

	"chat/internal/auth"
	"chat/internal/database"
	"chat/internal/models"
)

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Message struct {
	Data string `json:"data"`
}

var (
	messages  []map[string]interface{}
	clients   = make(map[*websocket.Conn]string)
	messageCh = make(chan Message)
	mu        sync.Mutex
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// регистрация
func HandlerSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userRequest UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	var user models.User
	result := database.DB.Where("username = ?", userRequest.Username).First(&user)
	if result.Error == nil {
		http.Error(w, "User already exist", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), 10)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	newUser := models.User{
		Username:     userRequest.Username,
		PasswordHash: string(hashedPassword),
	}
	result = database.DB.Create(&newUser)
	if result.Error != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	token := auth.JWTCreate(user.ID)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// авторизация
func HandlerSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userRequest UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	var user models.User
	result := database.DB.Where("username = ?", userRequest.Username).First(&user)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash), []byte(userRequest.Password),
	); err != nil {
		http.Error(w, "Wrong password", http.StatusBadRequest)
		return
	}

	token := auth.JWTCreate(user.ID)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func HandlerChat(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		log.Println("UserID not found in context")
		return
	}
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Println("User not found", err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrade error:", err)
		return
	}
	defer conn.Close()

	mu.Lock()
	clients[conn] = user.Username
	mu.Unlock()

	sendHistory(conn)

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			delete(clients, conn)
			break
		}
		fullMsg := map[string]interface{}{
			"username": clients[conn],
			"data":     msg.Data,
		}

		mu.Lock()
		messages = append(messages, fullMsg)
		mu.Unlock()

		messageCh <- msg
	}
}

func HandleMessages() {
	for {
		msg := <-messageCh

		mu.Lock()
		for client := range clients {
			fullMsg := map[string]interface{}{
				"username": clients[client],
				"data":     msg.Data,
			}

			if err := client.WriteJSON(fullMsg); err != nil {
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

// отравка истории сообщений подключившемуся пользователю
func sendHistory(conn *websocket.Conn) {
	for _, msg := range messages {
		if err := conn.WriteJSON(msg); err != nil {
			return
		}
	}
}
