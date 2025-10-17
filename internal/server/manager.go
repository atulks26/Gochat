package server

import (
	"log"
	"sync"
)

type ClientManager struct {
	clients     map[*User]bool
	clientsByID map[int64]*User
	mutex       sync.RWMutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients:     make(map[*User]bool),
		clientsByID: make(map[int64]*User),
	}
}

func (manager *ClientManager) AddClient(user *User) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	manager.clients[user] = true
	manager.clientsByID[user.ID] = user
	log.Printf("Added User %d to manager. Total clients: %d", user.ID, len(manager.clients))
}

func (manager *ClientManager) RemoveClient(user *User) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if manager.clients[user] {
		delete(manager.clients, user)
		delete(manager.clientsByID, user.ID)

		log.Printf("Removed User %d from manager. Total clients: %d", user.ID, len(manager.clients))
	}
}

func (manager *ClientManager) FindClientByID(userID int64) (*User, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	user, found := manager.clientsByID[userID]
	return user, found
}

//Broadcast function
