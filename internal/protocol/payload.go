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
	dst = EncodeString(dst, identity)
	dst = EncodeString(dst, password)

	return dst
}

func EncodeAuthSuccess(dst []byte, uid int64, username string) []byte {
	dst = binary.BigEndian.AppendUint64(dst, uint64(uid))
	dst = EncodeString(dst, username)
	return dst
}

func EncodeSendMessage(dst []byte, reciever, message string) []byte {
	dst = EncodeString(dst, reciever)
	dst = EncodeLongString(dst, message)

	return dst
}

func EncodeReceiveMessage(dst []byte, msg *users.MessageStored) []byte {
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

func EncodeChatListItem(dst []byte, partnerID int64, partnerUsername string, msg *users.MessageStored) []byte {
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

func DecodeChatListItem(b *bytes.Buffer) (*users.ChatPreview, error) {
	chat := &users.ChatPreview{
		Message: &users.MessageStored{},
	}

	var reader [8]byte

	if _, err := b.Read(reader[:]); err != nil {
		return nil, errors.New("error reading partner id")
	}
	chat.PartnerID = int64(binary.BigEndian.Uint64(reader[:]))

	username, err := DecodeString(b)
	if err != nil {
		return nil, err
	}
	chat.PartnerUsername = username

	if _, err := b.Read(reader[:]); err != nil {
		return nil, errors.New("error reading message id")
	}
	chat.Message.ID = int64(binary.BigEndian.Uint64(reader[:]))

	if _, err := b.Read(reader[:]); err != nil {
		return nil, errors.New("error reading sender id")
	}
	chat.Message.Sender_id = int64(binary.BigEndian.Uint64(reader[:]))

	if _, err := b.Read(reader[:]); err != nil {
		return nil, errors.New("error reading timestamp")
	}
	chat.Message.Timestamp = int64(binary.BigEndian.Uint64(reader[:]))

	flag, err := b.ReadByte()
	if err != nil {
		return nil, err
	}
	chat.Message.Is_read = (flag == 1)

	flag, err = b.ReadByte()
	if err != nil {
		return nil, err
	}
	chat.Message.Is_delivered = (flag == 1)

	content, err := DecodeLongString(b)
	if err != nil {
		return nil, err
	}
	chat.Message.Content = content

	return chat, nil
}

func DecodeChatSend(b *bytes.Buffer) (*users.MessageSent, error) {
	msg := &users.MessageSent{}

	var reader [8]byte

	if _, err := b.Read(reader[:]); err != nil {
		return nil, errors.New("error reading receiver id")
	}
	msg.Receiver_id = int64(binary.BigEndian.Uint64(reader[:]))

	content, err := DecodeLongString(b)
	if err != nil {
		return nil, err
	}

	msg.Content = content

	return msg, nil
}

func ParseAuth(b *bytes.Buffer) (string, string, error) {
	username, err := DecodeString(b)
	if err != nil {
		return "", "", err
	}

	password, err := DecodeString(b)
	if err != nil {
		return "", "", err
	}

	return username, password, nil
}

func DecodeAuthSuccess(b *bytes.Buffer) (int64, string, error) {
	if b.Len() < 9 {
		return 0, "", errors.New("payload too short")
	}

	uidBytes := make([]byte, 8)
	_, err := b.Read(uidBytes)
	if err != nil {
		return 0, "", err
	}

	uid := int64(binary.BigEndian.Uint64(uidBytes))

	username, err := DecodeString(b)
	if err != nil {
		return 0, "", err
	}

	return uid, username, nil
}
