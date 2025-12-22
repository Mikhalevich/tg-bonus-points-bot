package messagesender

import (
	"context"

	"github.com/go-telegram/bot"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

func (m *messageSender) DeleteMessage(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
) {
	if _, err := m.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    chatID.Int64(),
		MessageID: messageID.Int(),
	}); err != nil {
		logger.FromContext(ctx).
			WithError(err).
			Error("failed to delete message")
	}
}
