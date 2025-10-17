package server

import (
	"bufio"
	"chat/internal/helper"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type User struct {
	ID   int64
	Conn net.Conn
}

type Message struct {
	Source      int64
	Destination int64
	Mess        string
	TimeStamp   string
}

var nextUserID int64 = 0

func handleConnection(c net.Conn, manager *ClientManager) {
	userID := atomic.AddInt64(&nextUserID, 1)
	user := &User{
		ID:   userID,
		Conn: c,
	}

	manager.AddClient(user)
	defer manager.RemoveClient(user)

	log.Printf("New connection: User %d from %s", user.ID, user.Conn.RemoteAddr())
	defer log.Printf("User %d disconnected", user.ID)

	c.Write([]byte(fmt.Sprintf("Welcome, User %d\n", user.ID)))

	reader := bufio.NewReader(c)
	for {
		rawMessage, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Printf("Error reading from User %d: %v\n", user.ID, err)
			return
		}

		destID, messageStr, err := validateMessage(rawMessage)

		if err != nil {
			errResponse := fmt.Sprintf("%s\n", err.Error())
			c.Write([]byte(errResponse))
		} else {
			response, err := sendMessage(user.ID, destID, messageStr, manager)
			if err != nil {
				c.Write([]byte(err.Error()))
				continue
			}

			c.Write([]byte(response))
		}
	}
}

func validateMessage(message string) (int64, string, error) {
	trimmedMsg := strings.TrimSpace(message)
	parts := strings.SplitN(trimmedMsg, " ", 2)

	if len(parts) < 2 {
		return -1, "", errors.New("invalid format. Use <destination_id> <message>")
	}

	destIDStr := parts[0]
	messageStr := parts[1]

	destID, err := strconv.ParseInt(destIDStr, 10, 64)
	if err != nil {
		return -1, "", errors.New("invalid destination_id. It must be a number")
	}

	return destID, messageStr, nil
}

func sendMessage(srcID int64, destID int64, messageStr string, manager *ClientManager) (string, error) {
	message := &Message{
		Source:      srcID,
		Destination: destID,
		Mess:        messageStr,
		TimeStamp:   helper.FormatTime(time.Now()),
	}

	sendErr := MessageRouter(message, manager)
	if sendErr != nil {
		res := fmt.Sprintf("Message to User %d was not delivered. Reason: %v\n", message.Destination, sendErr)
		log.Printf("Message from User %d to User %d was not delivered. Reason: %v\n", message.Source, message.Destination, sendErr)

		return "", errors.New(res)
	}

	res := fmt.Sprintf("Message to User %d was successfully delivered.\n", message.Destination)
	log.Printf("Message from User %d to User %d at %v\n", message.Source, message.Destination, message.TimeStamp)

	return res, nil
}
