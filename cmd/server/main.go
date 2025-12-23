package main

import (
	"chat/internal/server"
	"chat/store/postgres"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	connStr := os.Getenv("DATABASE_URL")

	userStore, err := postgres.NewStore(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Starting chat server...")
	listener, err := server.StartServer(":8000")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	server.AcceptConnections(listener, server.NewOnlineClientManager(), userStore)
}
