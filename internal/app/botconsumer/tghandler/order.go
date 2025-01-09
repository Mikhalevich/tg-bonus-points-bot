package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/internal/tgbot"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

func (t *TGHandler) Order(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if err := t.orderProcessor.MakeOrder(
		ctx,
		msginfo.Info{
			ChatID:    msginfo.ChatIDFromInt(msg.ChatID),
			MessageID: msginfo.MessageIDFromInt(msg.MessageID),
		},
	); err != nil {
		return fmt.Errorf("make order: %w", err)
	}

	return nil
}
