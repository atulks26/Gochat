package server

import (
	"bufio"
	"chat/internal/auth"
	"chat/internal/helper"
	"chat/internal/protocol"
	"chat/store/users"
	"errors"
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

	var user *auth.User = nil

	for {
		frame, err := protocol.FrameRead(reader)
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Printf("Error reading frame: %v", err)
			break
		}

		if user == nil {
			userAuth, err := handleAuth(frame, userTable, manager)
			if err != nil {
				errMsg := protocol.EncodeLongString(err.Error())

				err := protocol.FrameWrite(c, protocol.OpError, errMsg)
				if err != nil {
					return
				}

				continue
			}

			user = userAuth

			manager.AddClient(user)
			defer manager.RemoveClient(user)

			log.Printf("New connection: %s", user.Username)
			defer log.Printf("%s disconnected", user.Username)

			authPayload := protocol.EncodeAuthSuccess(user.UID, user.Username)
			err2 := protocol.FrameWrite(c, protocol.OpAuthSuccess, authPayload)
			if err2 != nil {
				return
			}

			queue.ProcessOfflineMessages(user)
		} else {
			//send the message
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

func handleAuth(frame *protocol.Frame, userTable users.UserStore, manager Manager) (*auth.User, error) {
	opCode := frame.OpCode

	switch opCode {
	case protocol.OpRegister:
		return auth.ProcessLogin(frame.Payload, userTable, manager)
	case protocol.OpLogin:
		return auth.ProcessRegisteration(frame.Payload, userTable)
	default:
		return nil, errors.New("user not authenticated")
	}
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
