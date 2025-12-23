package protocol

import (
	"bytes"
	"chat/store/users"
	"encoding/binary"
	"errors"
)

// encode
func EncodeString(dst []byte, s string) []byte {
	if len(s) > 255 {
		s = s[:255]
	}

	dst = append(dst, byte(len(s)))
	dst = append(dst, s...)

	return dst
}

func EncodeLongString(dst []byte, s string) []byte {
	if len(s) > 65533 {
		s = s[:65533]
	}

	dst = binary.BigEndian.AppendUint16(dst, uint16(len(s)))
	dst = append(dst, s...)

	return dst
}

func EncodeAuth(dst []byte, identity, password string) []byte {
	return append(EncodeString(dst, identity), EncodeString(dst, password)...)
}

func EncodeAuthSuccess(dst []byte, uid int64, username string) []byte {
	dst = binary.BigEndian.AppendUint64(dst, uint64(uid))
	dst = EncodeString(dst, username)
	return dst
}

func EncodeSendMessage(dst []byte, reciever, message string) []byte {
	return append(EncodeString(dst, reciever), EncodeLongString(dst, message)...)
}

func EncodeReceiveMessage(dst []byte, msg *users.Message) []byte {
	dst = binary.BigEndian.AppendUint64(dst, uint64(msg.ID))
	dst = binary.BigEndian.AppendUint64(dst, uint64(msg.Sender_id))
	dst = binary.BigEndian.AppendUint64(dst, uint64(msg.Timestamp))

	if msg.Is_read {
		dst = append(dst, 1)
	} else {
		dst = append(dst, 0)
	}

	if msg.Is_delivered {
		dst = append(dst, 1)
	} else {
		dst = append(dst, 0)
	}

	dst = EncodeLongString(dst, msg.Content)
	return dst
}

func EncodeChatListItem(dst []byte, partnerID int64, partnerUsername string, msg *users.Message) []byte {
	dst = binary.BigEndian.AppendUint64(dst, uint64(partnerID))

	if len(partnerUsername) > 255 {
		partnerUsername = partnerUsername[:255]
	}
	dst = append(dst, byte(len(partnerUsername)))
	dst = append(dst, partnerUsername...)

	dst = EncodeReceiveMessage(dst, msg)

	return dst
}

// decode
func DecodeString(b *bytes.Buffer) (string, error) {
	lenByte, err := b.ReadByte()
	if err != nil {
		return "", errors.New("error reading string length")
	}

	lenB := int(lenByte)

	if b.Len() < lenB {
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

	if b.Len() < lenB {
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

func DecodeAuthSuccess(payload []byte) (int64, string, error) {
	if len(payload) < 9 {
		return 0, "", errors.New("payload too short")
	}

	buf := bytes.NewBuffer(payload)

	uidBytes := make([]byte, 8)
	_, err := buf.Read(uidBytes)
	if err != nil {
		return 0, "", err
	}

	uid := int64(binary.BigEndian.Uint64(uidBytes))

	username, err := DecodeString(buf)
	if err != nil {
		return 0, "", err
	}

	return uid, username, nil
}
