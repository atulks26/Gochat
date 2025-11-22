package protocol

import (
	"chat/internal/helper"
	"encoding/binary"
	"errors"
	"io"
	"net"
)

type Frame struct {
	OpCode  OpCode
	Payload []byte
}

var ErrPayloadTooLarge = errors.New("payload too large")

func FrameWrite(c net.Conn, op OpCode, payload []byte) error {
	pLen := len(payload)
	if pLen > 65535 {
		return ErrPayloadTooLarge
	}

	header := make([]byte, 3)
	header[0] = byte(op)
	binary.BigEndian.PutUint16(header[1:3], uint16(pLen))

	return helper.SafeWrite(c, append(header, payload...))
}

func FrameRead(reader io.Reader) (*Frame, error) {
	header := make([]byte, 3)

	if _, err := io.ReadFull(reader, header); err != nil {
		return nil, err
	}

	opCode := OpCode(header[0])
	pLen := binary.BigEndian.Uint16(header[1:3])

	payload := make([]byte, pLen)
	if _, err := io.ReadFull(reader, payload); err != nil {
		return nil, err
	}

	return &Frame{
		OpCode:  opCode,
		Payload: payload,
	}, nil
}
