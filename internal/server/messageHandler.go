package server

import (
	"chat/internal/protocol"
	"chat/store/users"
	"net"
	"sync"
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

	chat := func(partnerID int64, partnerUsername string, msg *users.Message) error {
		var payload []byte
		payload = protocol.EncodeChatListItem(payload, partnerID, partnerUsername, msg)

		return protocol.FrameWrite(c, protocol.OpChatListItem, payload)
	}

	if err := userTable.GetRecentChats(userID, chat); err != nil {
		return err
	}

	return protocol.FrameWrite(c, protocol.OpEndOfList, []byte{})
}
