package server

import (
	"bufio"
	"chat/internal/auth"
	"chat/internal/helper"
	"chat/store/users"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

type Manager interface {
	auth.OnlineUserChecker
	ClientFinder
	ClientLifecycleManager
}

type Message struct {
	Source      int64
	Sender      string
	Destination int64
	Mess        string
	TimeStamp   string
}

func handleConnection(c net.Conn, manager Manager, queue OfflineMessageQueue, userTable users.UserStore) {
	defer c.Close()
	reader := bufio.NewReader(c)

	user, close := auth.AuthenticateUser(c, reader, userTable, manager)
	if user == nil {
		return
	}

	if close {
		c.Close()
	}

	manager.AddClient(user)
	defer manager.RemoveClient(user)

	log.Printf("New connection: %s", user.Username)
	defer log.Printf("%s disconnected", user.Username)

	if err := helper.SafeWrite(c, []byte(fmt.Sprintf("Welcome, %s\n", user.Username))); err != nil {
		return
	}

	queue.ProcessOfflineMessages(user)

	for {
		rawMessage, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Printf("Error reading from User %d: %v\n", user.UID, err)
			c.Close()
			return
		}

		destUsername, messageStr, err := validateMessage(rawMessage)
		if err != nil {
			if err := helper.SafeWrite(c, []byte(fmt.Sprintf("%v\n", err.Error()))); err != nil {
				return
			}

			continue
		}

		destID, exists, err := userTable.FindClientByUsername(destUsername)
		if err != nil {
			// handle db error
		}

		if !exists {
			if err := helper.SafeWrite(c, []byte("User not found")); err != nil {
				return
			}
		}

		response, err := sendMessage(user.Username, user.UID, destID, messageStr, manager, queue)
		if err != nil {
			if err := helper.SafeWrite(c, []byte(err.Error())); err != nil {
				return
			}

			continue
		}

		if err := helper.SafeWrite(c, []byte(response)); err != nil {
			return
		}
	}
}

func validateMessage(message string) (string, string, error) {
	trimmedMsg := strings.TrimSpace(message)
	parts := strings.SplitN(trimmedMsg, " ", 2)

	if len(parts) < 2 {
		return "", "", errors.New("invalid format. Use <destination_username> <message>")
	}

	destUsername := parts[0]
	messageStr := parts[1]

	return destUsername, messageStr, nil
}

func sendMessage(sender string, srcID int64, destID int64, messageStr string, manager ClientFinder, queue OfflineMessageQueue) (string, error) {
	message := &Message{
		Source:      srcID,
		Sender:      sender,
		Destination: destID,
		Mess:        messageStr,
		TimeStamp:   helper.FormatTime(time.Now()),
	}

	res := MessageRouter(message, manager, queue)
	log.Printf("Message from User %d to User %d at %v\n", message.Source, message.Destination, message.TimeStamp)

	return res, nil
}
