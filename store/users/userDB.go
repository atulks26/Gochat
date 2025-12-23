package users

import (
	"time"
)

type UserData struct {
	ID             int64
	Username       string
	HashedPassword string
	CreatedAt      time.Time
}

type Message struct {
	ID           int64
	Sender_id    int64
	Timestamp    int64
	Is_read      bool
	Is_delivered bool
	Content      string
}

type RecentChat func(partnerID int64, partnerUsername string, msg *Message) error

type UserStore interface {
	FindClientByUsername(username string) (int64, bool, error)
	FindClientByID(userID int64) (*UserData, error)
	CreateUser(username string, hashedPassword string) (*UserData, error)
	GetRecentChats(userID int64, chat RecentChat) error
}
