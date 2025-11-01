package users

import (
	"sync"
	"sync/atomic"
	"time"
)

type UserData struct {
	ID             int64
	Username       string
	HashedPassword string
	CreatedAt      time.Time
}

type UserStore interface {
	FindClientByUsername(username string) (int64, bool, error)
	FindClientByID(userID int64) (*UserData, error)
	CreateUser(username string, hashedPassword string) (*UserData, error)
}

type UserTable struct {
	usersByID       map[int64]*UserData
	usersByUsername map[string]int64
	lastUID         int64
	mutex           sync.RWMutex
}

func NewUserTable() *UserTable {
	return &UserTable{
		usersByID:       make(map[int64]*UserData),
		usersByUsername: make(map[string]int64),
		lastUID:         000,
	}
}

func (userStore *UserTable) FindClientByUsername(username string) (int64, bool, error) {
	// handle db errors
	userStore.mutex.Lock()
	defer userStore.mutex.Unlock()

	id, found := userStore.usersByUsername[username]
	if !found {
		return -1, false, nil
	}

	return id, true, nil
}

func (userStore *UserTable) FindClientByID(userID int64) (*UserData, error) {
	// handle db errors
	userStore.mutex.Lock()
	defer userStore.mutex.Unlock()

	user := userStore.usersByID[userID]
	return user, nil
}

func (userStore *UserTable) CreateUser(username string, hashedPassword string) (*UserData, error) {
	// save to db, handle db errors
	userStore.mutex.Lock()
	defer userStore.mutex.Unlock()

	uid := atomic.AddInt64(&userStore.lastUID, 1)

	user := &UserData{
		ID:             uid,
		Username:       username,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now(),
	}

	userStore.usersByID[uid] = user
	userStore.usersByUsername[user.Username] = uid

	return user, nil
}
