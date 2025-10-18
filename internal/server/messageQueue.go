package server

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type MessageQueue struct {
	Buffer map[int64][]string
	mutex  sync.RWMutex
}

func NewMessageQueue() *MessageQueue {
	return &MessageQueue{
		Buffer: make(map[int64][]string),
	}
}

func (queue *MessageQueue) StoreOfflineMessage(message *Message, manager *ClientManager) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	_, ok := queue.Buffer[message.Destination]
	if !ok {
		queue.Buffer[message.Destination] = []string{}
	}

	offlineMessage := strconv.Itoa(int(message.Source)) + " " + message.Mess

	queue.Buffer[message.Destination] = append(queue.Buffer[message.Destination], offlineMessage)
}

func (queue *MessageQueue) ProcessOfflineMessages(user *User) {
	queue.mutex.Lock()

	messages, ok := queue.Buffer[user.ID]
	if !ok {
		queue.mutex.Unlock()
		return
	}

	delete(queue.Buffer, user.ID)
	queue.mutex.Unlock()

	for _, message := range messages {
		trimmedMsg := strings.TrimSpace(message)
		parts := strings.SplitN(trimmedMsg, " ", 2)

		srcIDStr := parts[0]
		messageStr := parts[1]

		// rep := fmt.Sprintf("(%v) User %d: %v\n", message.TimeStamp, srcIDStr, messageStr)
		rep := fmt.Sprintf("User %s: %s\n", srcIDStr, messageStr)

		user.Conn.Write([]byte(rep))
	}
}
