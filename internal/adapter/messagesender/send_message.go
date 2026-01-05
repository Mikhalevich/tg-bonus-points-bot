package messagesender

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

func (m *messageSender) SendText(
	ctx context.Context,
	chatID msginfo.ChatID,
	text string,
	rows ...button.InlineKeyboardButtonRow,
) {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID.Int64(),
		Text:        text,
		ReplyMarkup: makeButtonsMarkup(rows...),
	}); err != nil {
		logger.FromContext(ctx).
			WithError(err).
			WithField("text", text).
			Error("failed to send text")
	}
}

func (m *messageSender) SendTextMarkdown(
	ctx context.Context,
	chatID msginfo.ChatID,
	text string,
	rows ...button.InlineKeyboardButtonRow,
) {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID.Int64(),
		ParseMode:   models.ParseModeMarkdown,
		Text:        text,
		ReplyMarkup: makeButtonsMarkup(rows...),
	}); err != nil {
		logger.FromContext(ctx).
			WithError(err).
			WithField("text_markdown", text).
			Error("failed to send text markdown")
	}
}
