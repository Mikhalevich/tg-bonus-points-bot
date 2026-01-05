package messagesender

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

func (m *messageSender) ReplyText(
	ctx context.Context,
	chatID msginfo.ChatID,
	replyToMsgID msginfo.MessageID,
	text string,
	rows ...button.InlineKeyboardButtonRow,
) {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID.Int64(),
		ReplyParameters: &models.ReplyParameters{
			MessageID: replyToMsgID.Int(),
		},
		Text:        text,
		ReplyMarkup: makeButtonsMarkup(rows...),
	}); err != nil {
		logger.FromContext(ctx).
			WithError(err).
			WithField("text_plain", text).
			Error("failed to reply text")
	}
}

func (m *messageSender) ReplyTextMarkdown(
	ctx context.Context,
	chatID msginfo.ChatID,
	replyToMsgID msginfo.MessageID,
	text string,
	rows ...button.InlineKeyboardButtonRow,
) {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID.Int64(),
		ReplyParameters: &models.ReplyParameters{
			MessageID: replyToMsgID.Int(),
		},
		ParseMode:   models.ParseModeMarkdown,
		Text:        text,
		ReplyMarkup: makeButtonsMarkup(rows...),
	}); err != nil {
		logger.FromContext(ctx).
			WithError(err).
			WithField("text_markdown", text).
			Error("failed to reply text markdown")
	}
}
