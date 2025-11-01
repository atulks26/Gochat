package server

import (
	"chat/internal/auth"
	"log"
	"sync"
)

type ClientFinder interface {
	FindClientByID(userID int64) (*auth.User, bool)
}

type ClientLifecycleManager interface {
	AddClient(user *auth.User)
	RemoveClient(user *auth.User)
}

type OnlineClientManager struct {
	clientsByID map[int64]*auth.User
	mutex       sync.RWMutex
}

func NewOnlineClientManager() *OnlineClientManager {
	return &OnlineClientManager{
		clientsByID: make(map[int64]*auth.User),
	}
}

func (manager *OnlineClientManager) AddClient(user *auth.User) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	manager.clientsByID[user.UID] = user
	log.Printf("Added User %d to manager. Total clients: %d", user.UID, len(manager.clientsByID))
}

func (manager *OnlineClientManager) RemoveClient(user *auth.User) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	_, found := manager.clientsByID[user.UID]
	if found {
		delete(manager.clientsByID, user.UID)

		log.Printf("Removed User %d from manager. Total clients: %d", user.UID, len(manager.clientsByID))
	}
}

func (manager *OnlineClientManager) FindClientByID(userID int64) (*auth.User, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	user, found := manager.clientsByID[userID]
	return user, found
}

//Broadcast function
