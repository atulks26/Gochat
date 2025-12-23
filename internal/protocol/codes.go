package protocol

type OpCode uint8

const (
	OpRegister OpCode = 0x01
	OpLogin    OpCode = 0x02

	OpMessageSend    OpCode = 0x03
	OpMessageReceive OpCode = 0x04

	OpChatListItem   OpCode = 0x05
	OpEndOfList      OpCode = 0x06
	OpGetRecentChats OpCode = 0x07

	OpAuthSuccess OpCode = 0x10
	OpServerResp  OpCode = 0x11
	OpError       OpCode = 0xFF
)
