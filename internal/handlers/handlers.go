package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"chat/internal/config"
	"chat/internal/database"
	"chat/internal/models"
)

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Message struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

var (
	messages []Message
	mu       sync.Mutex
)

// регистрация
func HandleSignUp(w http.ResponseWriter, r *http.Request) {
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
	w.WriteHeader(http.StatusOK)
}

// авторизация
func HandleSignIn(w http.ResponseWriter, r *http.Request) {
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

	payload := jwt.MapClaims{
		"username": user.ID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString([]byte(config.AppConfig.JwtSecretKey))
	if err != nil {
		http.Error(w, "JWT Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": t,
	})
}

// ээ соси поменять на select
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		HandlePost(w, r)
	case http.MethodGet:
		HandleGet(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// обработка post запроса
func HandlePost(w http.ResponseWriter, r *http.Request) {
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if msg.Name == "" || msg.Data == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
	}

	mu.Lock()
	messages = append(messages, msg)
	mu.Unlock()
	json.NewEncoder(w).Encode(msg)
}

// обработка get запроса
func HandleGet(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	json.NewEncoder(w).Encode(messages)
}
