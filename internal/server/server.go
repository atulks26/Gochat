package server

import (
	"chat/store/users"
	"fmt"
	"log"
	"net"
)

func StartServer(addr string) (net.Listener, error) {
	l, err := net.Listen("tcp4", addr)
	if err != nil {
		return nil, err
	}

	fmt.Println("Server started on:", l.Addr().String())
	return l, nil
}

func AcceptConnections(l net.Listener, manager Manager, userTable users.UserStore) {
	defer l.Close()

	messages := NewMessageHandler()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Printf("Accept error: %v. Server shutting down", err)
			return
		}

		go handleConnection(c, manager, messages, userTable)
	}
}
