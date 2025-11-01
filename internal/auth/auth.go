package auth

import (
	"bufio"
	"chat/internal/helper"
	"chat/store/users"
	"io"
	"log"
	"net"
	"strings"
)

type User struct {
	UID      int64
	Username string
	Conn     net.Conn
}

type OnlineUserChecker interface {
	FindClientByID(userID int64) (*User, bool)
}

func (u *User) ID() int64 {
	return u.UID
}

func (u *User) Connection() net.Conn {
	return u.Conn
}

func AuthenticateUser(c net.Conn, reader *bufio.Reader, userTable users.UserStore, manager OnlineUserChecker) (*User, bool) {
	for {
		if err := helper.SafeWrite(c, []byte("Welcome! Login OR Signup? (L/R)\n")); err != nil {
			return nil, true
		}

		authOption, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil, false
			}

			log.Printf("Error reading authOption from User: %v\n", err)
			return nil, true
		}

		if authOption == "L\n" || authOption == "l\n" {
			user, terminate := processLogin(c, reader, userTable, manager)
			if user == nil && !terminate {
				continue
			}

			return user, terminate
		} else if authOption == "R\n" || authOption == "r\n" {
			user, terminate := processRegisteration(c, reader, userTable)
			if user == nil && !terminate {
				continue
			}

			return user, terminate
		} else {
			if err := helper.SafeWrite(c, []byte("Invalid choice\n")); err != nil {
				return nil, true
			}
		}
	}
}

func processLogin(c net.Conn, reader *bufio.Reader, userTable users.UserStore, manager OnlineUserChecker) (*User, bool) {
	if err := helper.SafeWrite(c, []byte("Enter username: ")); err != nil {
		return nil, true
	}

	usernameRaw, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return nil, false
		}

		log.Printf("Error reading Username from User: %v\n", err)
		return nil, true
	}
	username := strings.TrimRight(usernameRaw, "\r\n")

	id, isRegistered, err := userTable.FindClientByUsername(username)
	if err != nil {
		// handle db error
	}

	if isRegistered {
		if _, found := manager.FindClientByID(id); found {
			if err := helper.SafeWrite(c, []byte("This user is already logged in.\n")); err != nil {
				return nil, true
			}

			return nil, true
		}

		userData, err := userTable.FindClientByID(id)
		if err != nil {
			// handle db error
		}

		for i := 0; i < 3; i++ {
			if err := helper.SafeWrite(c, []byte("Enter password: ")); err != nil {
				return nil, true
			}

			passwordRaw, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					return nil, false
				}

				log.Printf("Error reading Password from User: %v\n", err)
				return nil, true
			}
			password := strings.TrimRight(passwordRaw, "\r\n")

			// hash password

			if password != userData.HashedPassword {
				if err := helper.SafeWrite(c, []byte("Wrong password\n")); err != nil {
					return nil, true
				}

				continue
			} else {
				user := &User{
					UID:      id,
					Username: username,
					Conn:     c,
				}

				if err := helper.SafeWrite(c, []byte("Logged in successfully\n")); err != nil {
					return nil, true
				}

				return user, false
			}
		}

		if err := helper.SafeWrite(c, []byte("Too many failed login attempts\n")); err != nil {
			return nil, true
		}

		return nil, true
	}

	if err := helper.SafeWrite(c, []byte("User not registered\n")); err != nil {
		return nil, true
	}

	return nil, false
}

func processRegisteration(c net.Conn, reader *bufio.Reader, userTable users.UserStore) (*User, bool) {
	for {
		if err := helper.SafeWrite(c, []byte("Choose a username: ")); err != nil {
			return nil, true
		}

		usernameRaw, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil, false
			}

			log.Printf("Error reading Username from User: %v\n", err)
			return nil, true
		}
		username := strings.TrimRight(usernameRaw, "\r\n")

		_, isRegistered, err := userTable.FindClientByUsername(username)
		if err != nil {
			// handle db error
		}

		if isRegistered {
			if err := helper.SafeWrite(c, []byte("Username taken\n")); err != nil {
				return nil, true
			}

			continue
		} else {
			for {
				if err := helper.SafeWrite(c, []byte("Choose a password: ")); err != nil {
					return nil, true
				}

				passwordRaw, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						return nil, false
					}

					log.Printf("Error reading Password from User: %v\n", err)
					return nil, true
				}
				password := strings.TrimRight(passwordRaw, "\r\n")

				if err := helper.SafeWrite(c, []byte("Confirm password: ")); err != nil {
					return nil, true
				}

				confRaw, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						return nil, false
					}

					log.Printf("Error reading Confirm Password from User: %v\n", err)
					return nil, true
				}
				confPassword := strings.TrimRight(confRaw, "\r\n")

				if password != confPassword {
					if err := helper.SafeWrite(c, []byte("Password and Confirm Password do not match\n")); err != nil {
						return nil, true
					}

					continue
				}

				// hash password here or after each input

				userdata, err := userTable.CreateUser(username, password)
				if err != nil {
					// handle db error
				}

				if userdata != nil {
					user := &User{
						UID:      userdata.ID,
						Username: userdata.Username,
						Conn:     c,
					}

					if err := helper.SafeWrite(c, []byte("Account created. Logged in successfully\n")); err != nil {
						return nil, true
					}

					return user, false
				}
			}
		}
	}
}
