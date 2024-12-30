package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/internal/tgbot"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

func (t *TGHandler) GetActiveOrder(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if err := t.orderProcessor.GetActiveOrder(
		ctx,
		msginfo.ChatIDFromInt(msg.ChatID),
		msginfo.MessageIDFromInt(msg.MessageID),
	); err != nil {
		return fmt.Errorf("make order: %w", err)
	}

	return nil
}
