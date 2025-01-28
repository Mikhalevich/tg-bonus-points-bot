package port

import (
	"context"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

type MessageSender interface {
	SendText(ctx context.Context, chatID msginfo.ChatID, text string)
	SendTextMarkdown(ctx context.Context, chatID msginfo.ChatID, text string)
	ReplyText(ctx context.Context, chatID msginfo.ChatID, replyToMsgID msginfo.MessageID, text string,
		buttons ...button.InlineKeyboardButtonRow)
	ReplyTextMarkdown(ctx context.Context, chatID msginfo.ChatID,
		replyToMsgID msginfo.MessageID, text string, buttons ...button.InlineKeyboardButtonRow)
	EscapeMarkdown(s string) string
	SendPNGMarkdown(ctx context.Context, chatID msginfo.ChatID, caption string, png []byte,
		buttons ...button.InlineKeyboardButtonRow) error
	EditTextMessage(ctx context.Context, chatID msginfo.ChatID, messageID msginfo.MessageID,
		text string, rows ...button.InlineKeyboardButtonRow,
	)
	DeleteMessage(ctx context.Context, chatID msginfo.ChatID, messageID msginfo.MessageID)
	SendOrderInvoice(
		ctx context.Context,
		chatID msginfo.ChatID,
		title string,
		description string,
		ord order.Order,
		rows ...button.InlineKeyboardButtonRow,
	) error
	AnswerOrderPayment(ctx context.Context, paymentID string, ok bool, errorMsg string) error
}
