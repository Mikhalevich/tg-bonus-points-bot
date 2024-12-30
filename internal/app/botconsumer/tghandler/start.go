package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/internal/tgbot"
)

func (t *TGHandler) Start(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if err := sender.SendMessage(
		ctx,
		msg.ChatID,
		"Type /order for requesting an order",
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
