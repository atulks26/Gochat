package server

import (
	"fmt"
)

func MessageRouter(message *Message, manager *ClientManager, queue *MessageQueue) string {
	targetUser, found := manager.FindClientByID(message.Destination)
	if !found {
		queue.StoreOfflineMessage(message, manager)

		return "Message sent but not delivered yet.\n"
	}

	res := fmt.Sprintf("Message to User %d was successfully delivered.\n", message.Destination)

	rep := fmt.Sprintf("(%v) User %d: %v\n", message.TimeStamp, message.Source, message.Mess)
	targetUser.Conn.Write([]byte(rep))

	return res
}
