package server

import (
	"bufio"
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
	TimeStamp   time.Time
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

			log.Printf("Error reading from User %d: %v", user.ID, err)
			return
		}

		trimmedMsg := strings.TrimSpace(rawMessage)
		parts := strings.SplitN(trimmedMsg, " ", 2)

		if len(parts) < 2 {
			c.Write([]byte("Invalid format. Use <destination_id> <message>\n"))
			continue
		}

		destIDStr := parts[0]
		messageStr := parts[1]

		destID, err := strconv.ParseInt(destIDStr, 10, 64)
		if err != nil {
			c.Write([]byte("Invalid destination_id. It must be a number.\n"))
			continue
		}

		message := &Message{
			Source:      user.ID,
			Destination: destID,
			Mess:        messageStr,
		}

		log.Printf("Message from User %d to User %d: %s", message.Source, message.Destination, message.Mess)
		c.Write([]byte("Echo: " + message.Mess))
	}
}
