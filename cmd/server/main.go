package main

import (
	"chat/internal/server"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("please provide a port number")
	}

	port := ":" + os.Args[1]

	listener, err := server.StartServer(port)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}

	fmt.Println("Server initialized")

	manager := server.NewOnlineClientManager()

	server.AcceptConnections(listener, manager)
}
