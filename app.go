package main

import (
	"bytes"
	"chat/internal/protocol"
	"context"
	"errors"
	"fmt"
	"net"
)

type App struct {
	ctx context.Context
	c   net.Conn
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

func (a *App) Login(username string, password string) (string, error) {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		return "", errors.New("error establishing connection")
	}

	a.c = conn

	payload := protocol.EncodeAuth(username, password)
	err = protocol.FrameWrite(conn, protocol.OpLogin, payload)
	if err != nil {
		return "", err
	}

	frame, err := protocol.FrameRead(conn)
	if err != nil {
		return "", err
	}

	switch frame.OpCode {
	case protocol.OpAuthSuccess:
		return "Login Successful", nil
	case protocol.OpError:
		errMsg, _ := protocol.DecodeLongString(bytes.NewBuffer(frame.Payload))
		return "", errors.New(errMsg)
	default:
		return "", errors.New("unexpected server response")
	}
}

func (a *App) Register(email string, password string) (string, error) {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		return "", errors.New("error establishing connection")
	}

	a.c = conn

	payload := protocol.EncodeAuth(email, password)
	err = protocol.FrameWrite(conn, protocol.OpRegister, payload)
	if err != nil {
		return "", err
	}

	frame, err := protocol.FrameRead(conn)
	if err != nil {
		return "", err
	}

	switch frame.OpCode {
	case protocol.OpAuthSuccess:
		return "Registered successfully", nil
	case protocol.OpError:
		errMsg, _ := protocol.DecodeLongString(bytes.NewBuffer(frame.Payload))
		return "", errors.New(errMsg)
	default:
		return "", errors.New("unexpected server response")
	}
}
