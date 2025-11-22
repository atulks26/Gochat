package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// encode
func EncodeString(s string) []byte {
	if len(s) > 255 {
		s = s[:255]
	}

	b := make([]byte, len(s)+1)
	b[0] = byte(len(s))
	copy(b[1:], s)

	return b
}

func EncodeLongString(s string) []byte {
	if len(s) > 65533 {
		s = s[:65533]
	}

	b := make([]byte, len(s)+2)
	binary.BigEndian.PutUint16(b[0:2], uint16(len(s)))
	copy(b[2:], s)

	return b
}

func EncodeAuth(username, password string) []byte {
	return append(EncodeString(username), EncodeString(password)...)
}

func EncodeAuthSuccess(uid int64, username string) []byte {
	usernameByte := EncodeString(username)

	b := make([]byte, len(usernameByte)+8)
	binary.BigEndian.PutUint64(b[:8], uint64(uid))
	copy(b[8:], usernameByte)

	return b
}

func EncodeOutMessage(reciever, message string) []byte {
	return append(EncodeString(reciever), EncodeLongString(message)...)
}

func EncodeInMessage(timestamp, sender, message string) []byte {
	payload := append(EncodeString(timestamp), EncodeString(sender)...)
	return append(payload, EncodeLongString(message)...)
}

// decode
func DecodeString(b *bytes.Buffer) (string, error) {
	lenByte, err := b.ReadByte()
	if err != nil {
		return "", errors.New("error reading string length")
	}

	lenB := int(lenByte)

	if b.Len()-1 != lenB {
		return "", errors.New("incorrect string length")
	}

	return string(b.Next(lenB)), nil
}

func DecodeLongString(b *bytes.Buffer) (string, error) {
	lenBytes := make([]byte, 2)

	if _, err := b.Read(lenBytes); err != nil {
		return "", errors.New("error reading string length")
	}

	lenB := int(binary.BigEndian.Uint16(lenBytes))

	if b.Len()-2 != lenB {
		return "", errors.New("incorrect string length")
	}

	return string(b.Next(lenB)), nil
}

// make separate parser for registeration later
func ParseAuth(payload []byte) (string, string, error) {
	buf := bytes.NewBuffer(payload)

	username, err := DecodeString(buf)
	if err != nil {
		return "", "", err
	}

	password, err := DecodeString(buf)
	if err != nil {
		return "", "", err
	}

	return username, password, nil
}
