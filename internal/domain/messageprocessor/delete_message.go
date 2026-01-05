package messageprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (m *MessageProcessor) DeleteMessage(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
) error {
	if err := m.sender.DeleteMessage(ctx, chatID, messageID); err != nil {
		return fmt.Errorf("delete message: %w", err)
	}

	return nil
}
