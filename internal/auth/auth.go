package auth

import (
	"chat/internal/protocol"
	"chat/store/users"
	"errors"
	"net"
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
	username, password, err := protocol.ParseAuth(payload)
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

		//hash password, then check
		if password != userdata.HashedPassword {
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
	username, password, err := protocol.ParseAuth(payload)
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

	// add password hashing here, or on client side
	userdata, err := userTable.CreateUser(username, password)
	if err != nil {
		return nil, err
	}

	user := &User{
		UID:      userdata.ID,
		Username: userdata.Username,
	}

	return user, nil
}
