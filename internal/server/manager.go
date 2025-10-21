package server

import (
	"log"
	"sync"
)

type OnlineClientManager struct {
	clientsByID map[int64]*User
	mutex       sync.RWMutex
}

func NewOnlineClientManager() *OnlineClientManager {
	return &OnlineClientManager{
		clientsByID: make(map[int64]*User),
	}
}

func (manager *OnlineClientManager) AddClient(user *User) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	// find if user is already registered

	// yes: go to enter password

	// no: go to register

	manager.clientsByID[user.ID] = user
	log.Printf("Added User %d to manager. Total clients: %d", user.ID, len(manager.clientsByID))
}

func (manager *OnlineClientManager) RemoveClient(user *User) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	_, found := manager.FindClientByID(user.ID)
	if found {
		delete(manager.clientsByID, user.ID)

		log.Printf("Removed User %d from manager. Total clients: %d", user.ID, len(manager.clientsByID))
	}
}

func (manager *OnlineClientManager) FindClientByID(userID int64) (*User, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	user, found := manager.clientsByID[userID]
	return user, found
}

//Broadcast function
