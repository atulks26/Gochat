package helper

import (
	"log"
	"net"
)

func SafeWrite(c net.Conn, data []byte) error {
	_, err := c.Write(data)
	if err != nil {
		log.Printf("Error writing to user %s: %v", c.RemoteAddr(), err)
		return err
	}

	return nil
}
