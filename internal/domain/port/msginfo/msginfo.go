package msginfo

import (
	"context"
)

type MessageID int

func (m MessageID) Int() int {
	return int(m)
}

func MessageIDFromInt(id int) MessageID {
	return MessageID(id)
}

type ChatID int64

func (c ChatID) Int64() int64 {
	return int64(c)
}

func ChatIDFromInt(id int64) ChatID {
	return ChatID(id)
}

type Info struct {
	MessageID MessageID
	ChatID    ChatID
	// for text message
	Text string
	// for callback query
	Data string
}

type Handler func(ctx context.Context, info Info) error
