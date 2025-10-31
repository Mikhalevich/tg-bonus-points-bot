package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

func (t *TGHandler) Order(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if err := t.cartProcessor.Create(
		ctx,
		msginfo.Info{
			ChatID:    msginfo.ChatIDFromInt(msg.ChatID),
			MessageID: msginfo.MessageIDFromInt(msg.MessageID),
		},
	); err != nil {
		return fmt.Errorf("start new cart: %w", err)
	}

	return nil
}
