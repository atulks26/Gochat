package auth

import (
	"bytes"
	"chat/internal/protocol"
	"chat/store/users"
	"errors"
	"net"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UID      int64
	Username string
	Conn     net.Conn
}

type OnlineUserChecker interface {
	FindOnlineClientByID(userID int64) (*User, bool)
}

func (u *User) ID() int64 {
	return u.UID
}

func (u *User) Connection() net.Conn {
	return u.Conn
}

func ProcessLogin(payload []byte, userTable users.UserStore, manager OnlineUserChecker) (*User, error) {
	username, password, err := protocol.ParseAuth(bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	id, isRegistered, err := userTable.FindClientByUsername(username)
	if err != nil {
		return nil, err
	}

	if isRegistered {
		if _, found := manager.FindOnlineClientByID(id); found {
			return nil, errors.New("user is already logged in")
		}

		userdata, err := userTable.FindClientByID(id)
		if err != nil {
			return nil, err
		}

		err = bcrypt.CompareHashAndPassword([]byte(userdata.HashedPassword), []byte(password))
		if err != nil {
			return nil, errors.New("password incorrect")
		}

		user := &User{
			UID:      userdata.ID,
			Username: userdata.Username,
		}

		return user, nil
	}

	return nil, errors.New("user not found")
}

func ProcessRegisteration(payload []byte, userTable users.UserStore) (*User, error) {
	username, password, err := protocol.ParseAuth(bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	_, isRegistered, err := userTable.FindClientByUsername(username)
	if err != nil {
		return nil, err
	}

	if isRegistered {
		return nil, errors.New("username taken")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	userdata, err := userTable.CreateUser(username, string(passwordHash))
	if err != nil {
		return nil, err
	}

	user := &User{
		UID:      userdata.ID,
		Username: userdata.Username,
	}

	return user, nil
}
