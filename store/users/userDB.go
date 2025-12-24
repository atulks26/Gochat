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

type MessageStored struct {
	ID           int64
	Sender_id    int64
	Timestamp    int64
	Is_read      bool
	Is_delivered bool
	Content      string
}

type MessageSent struct {
	Sender_id   int64
	Receiver_id int64
	Timestamp   int64
	Content     string
}

type ChatPreview struct {
	PartnerID       int64          `json:"user_id"`
	PartnerUsername string         `json:"username"`
	Message         *MessageStored `json:"message"`
}

type RecentChat func(partnerID int64, partnerUsername string, msg *MessageStored) error

type UserStore interface {
	FindClientByUsername(username string) (int64, bool, error)
	FindClientByID(userID int64) (*UserData, error)
	CreateUser(username string, hashedPassword string) (*UserData, error)
	GetRecentChats(userID int64, chat RecentChat) error
	SaveNewMessage(msg *MessageSent) (int64, error)
}
