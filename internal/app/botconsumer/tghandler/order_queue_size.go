package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/internal/tgbot"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

func (t *TGHandler) OrderQueueSize(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if err := t.orderProcessor.OrderQueueSize(
		ctx,
		msginfo.Info{
			ChatID:    msginfo.ChatIDFromInt(msg.ChatID),
			MessageID: msginfo.MessageIDFromInt(msg.MessageID),
		},
	); err != nil {
		return fmt.Errorf("order queue size: %w", err)
	}

	return nil
}
