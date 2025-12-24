package main

import (
	"bytes"
	"chat/internal/protocol"
	"chat/store/users"
	"context"
	"errors"
	"fmt"
	"net"
)

type App struct {
	ctx context.Context
	c   net.Conn
}

type AuthUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) Login(username string, password string) (*AuthUser, error) {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		return nil, errors.New("error establishing connection")
	}

	a.c = conn

	var payload []byte

	payload = protocol.EncodeAuth(payload, username, password)
	err = protocol.FrameWrite(conn, protocol.OpLogin, payload)
	if err != nil {
		return nil, err
	}

	frame, err := protocol.FrameRead(conn)
	if err != nil {
		return nil, err
	}

	switch frame.OpCode {
	case protocol.OpAuthSuccess:
		payloadBuf := bytes.NewBuffer(frame.Payload)
		uid, name, err := protocol.DecodeAuthSuccess(payloadBuf)
		if err != nil {
			return nil, errors.New("failed to decode success payload")
		}

		return &AuthUser{ID: uid, Username: name}, nil
	case protocol.OpError:
		errMsg, _ := protocol.DecodeLongString(bytes.NewBuffer(frame.Payload))
		return nil, errors.New(errMsg)
	default:
		return nil, errors.New("unexpected server response")
	}
}

func (a *App) Register(email string, password string) (*AuthUser, error) {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		return nil, errors.New("error establishing connection")
	}

	a.c = conn

	var payload []byte
	payload = protocol.EncodeAuth(payload, email, password)
	err = protocol.FrameWrite(conn, protocol.OpRegister, payload)
	if err != nil {
		return nil, err
	}

	frame, err := protocol.FrameRead(conn)
	if err != nil {
		return nil, err
	}

	switch frame.OpCode {
	case protocol.OpAuthSuccess:
		uid, name, err := protocol.DecodeAuthSuccess(bytes.NewBuffer(frame.Payload))
		if err != nil {
			return nil, errors.New("failed to decode success payload")
		}

		return &AuthUser{ID: uid, Username: name}, nil
	case protocol.OpError:
		errMsg, _ := protocol.DecodeLongString(bytes.NewBuffer(frame.Payload))
		return nil, errors.New(errMsg)
	default:
		return nil, errors.New("unexpected server response")
	}
}

func (a *App) FetchRecentChats() ([]*users.ChatPreview, error) {
	if a.c == nil {
		return nil, errors.New("not connected")
	}

	err := protocol.FrameWrite(a.c, protocol.OpGetRecentChats, []byte{})
	if err != nil {
		return nil, err
	}

	var chats []*users.ChatPreview

	for {
		frame, err := protocol.FrameRead(a.c)
		if err != nil {
			return nil, err
		}

		switch frame.OpCode {
		case protocol.OpChatListItem:
			chat, err := protocol.DecodeChatListItem(bytes.NewBuffer(frame.Payload))
			if err != nil {
				return nil, err
			}

			chats = append(chats, chat)

		case protocol.OpEndOfList:
			return chats, nil

		case protocol.OpError:
			errMsg, _ := protocol.DecodeLongString(bytes.NewBuffer(frame.Payload))
			return nil, errors.New(errMsg)

		default:
			return nil, fmt.Errorf("unexpected opcode: %d", frame.OpCode)
		}
	}
}
