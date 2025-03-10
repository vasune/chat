package main

import (
	"chat/internal/database"
	"chat/internal/handlers"
	"log"
	"net/http"
)

func main() {
	database.SetupDB()

	http.HandleFunc("/chat", handlers.HandleConnections)
	http.HandleFunc("/signup", handlers.HandleSignUp)
	http.HandleFunc("/signin", handlers.HandleSignIn)

	log.Println("Server started")
	log.Println("negr")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
