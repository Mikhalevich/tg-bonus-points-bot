package port

import (
	"context"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

type Button struct {
	Text string
	Data string
}

type MessageSender interface {
	SendText(ctx context.Context, chatID msginfo.ChatID, text string)
	SendTextMarkdown(ctx context.Context, chatID msginfo.ChatID, text string)
	ReplyText(ctx context.Context, chatID msginfo.ChatID, replyToMsgID msginfo.MessageID, text string, buttons ...Button)
	ReplyTextMarkdown(ctx context.Context, chatID msginfo.ChatID,
		replyToMsgID msginfo.MessageID, text string, buttons ...Button)
	EscapeMarkdown(s string) string
	SendPNGMarkdown(ctx context.Context, chatID msginfo.ChatID, caption string, png []byte, buttons ...Button) error
}
