package port

import (
	"context"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

type MessageSender interface {
	SendText(ctx context.Context, chatID msginfo.ChatID, text string)
	SendTextMarkdown(ctx context.Context, chatID msginfo.ChatID, text string)
	ReplyText(ctx context.Context, chatID msginfo.ChatID, replyToMsgID msginfo.MessageID, text string,
		buttons ...button.InlineKeyboardButton)
	ReplyTextMarkdown(ctx context.Context, chatID msginfo.ChatID,
		replyToMsgID msginfo.MessageID, text string, buttons ...button.InlineKeyboardButton)
	EscapeMarkdown(s string) string
	SendPNGMarkdown(ctx context.Context, chatID msginfo.ChatID, caption string, png []byte,
		buttons ...button.InlineKeyboardButton) error
}
