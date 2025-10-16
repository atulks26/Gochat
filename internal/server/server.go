package server

import (
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

func AcceptConnections(l net.Listener, manager *ClientManager) {
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Printf("Accept error: %v. Server shutting down", err)
			return
		}

		go handleConnection(c, manager)
	}
}
