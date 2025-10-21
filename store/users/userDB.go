package users

import (
	"sync"
	"time"
)

type UserData struct {
	Username       string
	HashedPassword string
	CreatedAt      time.Time
}

type UserTable struct {
	usersByID       map[int64]*UserData
	usersByUsername map[string]int64
	mutex           sync.RWMutex
}

func NewUserTable() *UserTable {
	return &UserTable{
		usersByID:       make(map[int64]*UserData),
		usersByUsername: make(map[string]int64),
	}
}

func (userStore *UserTable) FindClientByUsername(username string) (int64, bool, error) {
	// handle error during actual DB search
	userStore.mutex.Lock()
	defer userStore.mutex.Unlock()

	id, found := userStore.usersByUsername[username]
	if !found {
		return -1, false, nil
	}

	return id, true, nil
}

func (userStore *UserTable) AddUser() (bool, error) {
	userStore.mutex.Lock()
	defer userStore.mutex.Unlock()

	// receive user data and add to table

	return true, nil
}
