package server

import (
	"log"
	"sync"
)

type ClientManager struct {
	clients map[*User]bool
	mutex   sync.RWMutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[*User]bool),
	}
}

func (manager *ClientManager) AddClient(user *User) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	manager.clients[user] = true
	log.Printf("Added User %d to manager. Total clients: %d", user.ID, len(manager.clients))
}

func (manager *ClientManager) RemoveClient(user *User) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if manager.clients[user] {
		delete(manager.clients, user)

		log.Printf("Removed User %d from manager. Total clients: %d", user.ID, len(manager.clients))
	}
}

//Broadcast function
