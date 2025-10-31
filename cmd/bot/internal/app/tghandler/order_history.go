package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

func (t *TGHandler) OrderHistory(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if err := t.historyProcessor.Show(
		ctx,
		msginfo.ChatIDFromInt(msg.ChatID),
	); err != nil {
		return fmt.Errorf("history orders: %w", err)
	}

	return nil
}
