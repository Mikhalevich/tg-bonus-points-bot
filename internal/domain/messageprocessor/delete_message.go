package messageprocessor

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (m *MessageProcessor) DeleteMessage(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
) {
	m.sender.DeleteMessage(ctx, chatID, messageID)
}
