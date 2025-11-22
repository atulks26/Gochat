package server

import (
	"fmt"
)

func MessageRouter(message *Message, manager ClientFinder, queue OfflineMessageQueue) string {
	targetUser, found := manager.FindOnlineClientByID(message.Destination)
	if !found {
		queue.StoreOfflineMessage(message)

		return "Message sent but not delivered yet.\n"
	}

	res := fmt.Sprintf("Message to %s was successfully delivered.\n", targetUser.Username)

	rep := fmt.Sprintf("(%v) %s: %v\n", message.TimeStamp, message.Sender, message.Mess)
	targetUser.Conn.Write([]byte(rep))

	return res
}
