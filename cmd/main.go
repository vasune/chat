package main

import (
	"chat/internal/auth"
	"chat/internal/config"
	"chat/internal/database"
	"chat/internal/handlers"
	"chat/internal/repository/postgres"
	"chat/internal/usecases"
	"log"

	"net/http"
)

func main() {
	config.LoadCfg()
	database.SetupDB()

	userRepo := postgres.NewUserRepository(database.DB)
	authUseCase := usecases.NewAuthUseCase(userRepo)
	authHandler := handlers.NewAuthHandler(authUseCase)

	chatUseCase := usecases.NewChatUseCase(userRepo)
	chatHandler := handlers.NewChatHandler(chatUseCase)

	go chatHandler.HandleMessages()

	http.Handle("/chat", auth.JWTMiddleware(chatHandler.HandleConnections))
	http.HandleFunc("/signup", authHandler.SignUp)
	http.HandleFunc("/signin", authHandler.SignIn)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
