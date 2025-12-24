package server

import (
	"bufio"
	"chat/internal/auth"
	"chat/internal/protocol"
	"chat/store/users"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
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

func handleConnection(c net.Conn, manager Manager, messages *MessageHandler, userTable users.UserStore) {
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
				var errMsg []byte
				errMsg = protocol.EncodeLongString(errMsg, err.Error())

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

			var authPayload []byte
			authPayload = protocol.EncodeAuthSuccess(authPayload, user.UID, user.Username)
			err2 := protocol.FrameWrite(c, protocol.OpAuthSuccess, authPayload)
			if err2 != nil {
				return
			}

		} else {
			switch frame.OpCode {
			case protocol.OpGetRecentChats:
				err := messages.FetchRecentChats(c, user.UID, userTable)
				if err != nil {
					if sendErr := SendError(c, err); sendErr != nil {
						fmt.Printf("Broken network: %v\n", sendErr)
						return
					}
					continue
				}

			case protocol.OpMessageSend:
				err := SendMessage(c, user.UID, *frame, userTable)
				if err != nil {
					if sendErr := SendError(c, err); sendErr != nil {
						fmt.Printf("Broken network: %v\n", sendErr)
						return
					}
					continue
				}

			default:
				if sendErr := SendError(c, err); sendErr != nil {
					fmt.Printf("Broken network: %v\n", sendErr)
					return
				}
			}
		}
	}
}

func handleAuth(frame *protocol.Frame, userTable users.UserStore, manager Manager) (*auth.User, error) {
	opCode := frame.OpCode

	switch opCode {
	case protocol.OpRegister:
		return auth.ProcessRegisteration(frame.Payload, userTable)
	case protocol.OpLogin:
		return auth.ProcessLogin(frame.Payload, userTable, manager)
	default:
		return nil, errors.New("user not authenticated")
	}
}

func SendError(c net.Conn, err error) error {
	payload := protocol.EncodeLongString(nil, err.Error())

	return protocol.FrameWrite(c, protocol.OpError, payload)
}
