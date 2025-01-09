package messagesender

import (
	"context"

	"github.com/go-telegram/bot"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/logger"
)

func (m *messageSender) EditTextMessage(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	text string,
	rows ...button.InlineKeyboardButtonRow,
) {
	if _, err := m.bot.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatID.Int64(),
		MessageID:   messageID.Int(),
		Text:        text,
		ReplyMarkup: makeButtonsMarkup(rows...),
	}); err != nil {
		logger.FromContext(ctx).
			WithError(err).
			WithField("text_plain", text).
			Error("failed to edit text")
	}
}
