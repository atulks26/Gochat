package server

import (
	"errors"
	"fmt"
)

func MessageRouter(message *Message, manager *ClientManager) (err error) {
	targetUser, found := manager.FindClientByID(message.Destination)
	if !found {
		return errors.New("User not found")
	}

	rep := fmt.Sprintf("(%v) User %d: %v\n", message.TimeStamp, message.Source, message.Mess)
	targetUser.Conn.Write([]byte(rep))

	return nil
}
