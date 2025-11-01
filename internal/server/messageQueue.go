package server

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

type MessageReceiver interface {
	ID() int64
	Connection() net.Conn
}

type OfflineMessageQueue interface {
	StoreOfflineMessage(message *Message)
	ProcessOfflineMessages(user MessageReceiver)
}

type MessageQueue struct {
	Buffer map[int64][]string
	mutex  sync.RWMutex
}

func NewMessageQueue() *MessageQueue {
	return &MessageQueue{
		Buffer: make(map[int64][]string),
	}
}

func (queue *MessageQueue) StoreOfflineMessage(message *Message) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	_, ok := queue.Buffer[message.Destination]
	if !ok {
		queue.Buffer[message.Destination] = []string{}
	}

	offlineMessage := message.Sender + " " + message.Mess

	queue.Buffer[message.Destination] = append(queue.Buffer[message.Destination], offlineMessage)
}

func (queue *MessageQueue) ProcessOfflineMessages(user MessageReceiver) {
	queue.mutex.Lock()

	messages, ok := queue.Buffer[user.ID()]
	if !ok {
		queue.mutex.Unlock()
		return
	}

	delete(queue.Buffer, user.ID())
	queue.mutex.Unlock()

	for _, message := range messages {
		trimmedMsg := strings.TrimSpace(message)
		parts := strings.SplitN(trimmedMsg, " ", 2)

		srcIDStr := parts[0]
		messageStr := parts[1]

		// rep := fmt.Sprintf("(%v) User %d: %v\n", message.TimeStamp, srcIDStr, messageStr)
		rep := fmt.Sprintf("%s: %s\n", srcIDStr, messageStr)

		user.Connection().Write([]byte(rep))
	}
}
