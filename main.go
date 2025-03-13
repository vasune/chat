package main

import (
	"chat/internal/auth"
	"chat/internal/database"
	"chat/internal/handlers"
	"log"
	"net/http"
)

func main() {
	database.SetupDB()

	go handlers.HandleMessages()

	http.Handle("/chat", auth.JWTMiddleware(http.HandlerFunc(handlers.HandlerChat)))

	http.HandleFunc("/signup", handlers.HandlerSignUp)
	http.HandleFunc("/signin", handlers.HandlerSignIn)

	log.Println("Server started")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
