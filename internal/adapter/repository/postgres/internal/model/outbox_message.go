package model

import (
	"database/sql"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/jsonb"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
)

type MessageType string

const (
	MessageTypePlain    MessageType = "plain"
	MessageTypeMarkdown MessageType = "markdown"
	MessageTypePNG      MessageType = "png"
)

type OutboxMessage struct {
	ID             int           `db:"id"`
	ChatID         int64         `db:"chat_id"`
	ReplyMessageID sql.NullInt64 `db:"reply_msg_id"`
	Text           string        `db:"msg_text"`
	Type           MessageType   `db:"msg_type"`
	Payload        []byte        `db:"payload"`
	Button         jsonb.JSONB   `db:"buttons"`
}

func ToDBMessageType(msgType messageprocessor.MessageTextType) MessageType {
	switch msgType {
	case messageprocessor.MessageTextTypePlain:
		return MessageTypePlain
	case messageprocessor.MessageTextTypeMarkdown:
		return MessageTypeMarkdown
	}

	return ""
}
