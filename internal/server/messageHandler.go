package server

import (
	"bytes"
	"chat/internal/protocol"
	"chat/store/users"
	"encoding/binary"
	"net"
	"sync"
	"time"
)

type IncomingMessageQueue interface {
	FetchRecentChats(userID int64)
}

type MessageHandler struct {
	MessageWaitQueue []string
	isFetching       bool
	mutex            sync.Mutex
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		MessageWaitQueue: make([]string, 0),
		isFetching:       false,
	}
}

func (h *MessageHandler) FetchRecentChats(c net.Conn, userID int64, userTable users.UserStore) error {
	h.mutex.Lock()
	h.isFetching = true
	h.mutex.Unlock()

	defer func() {
		h.mutex.Lock()
		h.isFetching = false
		h.mutex.Unlock()
	}()

	chat := func(partnerID int64, partnerUsername string, msg *users.MessageStored) error {
		var payload []byte
		payload = protocol.EncodeChatListItem(payload, partnerID, partnerUsername, msg)

		return protocol.FrameWrite(c, protocol.OpChatListItem, payload)
	}

	if err := userTable.GetRecentChats(userID, chat); err != nil {
		return err
	}

	return protocol.FrameWrite(c, protocol.OpEndOfList, []byte{})
}

func SendMessage(c net.Conn, uid int64, f protocol.Frame, userTable users.UserStore) error {
	msg, err := protocol.DecodeChatSend(bytes.NewBuffer(f.Payload))
	if err != nil {
		return err
	}

	msg.Sender_id = uid
	msg.Timestamp = time.Now().Unix()

	msgID, err := userTable.SaveNewMessage(msg)
	if err != nil {
		return err
	}

	var resp [16]byte
	binary.BigEndian.PutUint64(resp[0:8], uint64(msgID))
	binary.BigEndian.PutUint64(resp[8:16], uint64(msg.Timestamp))

	if err := protocol.FrameWrite(c, protocol.OpServerResp, resp[:]); err != nil {
		return err
	}

	return nil
}
