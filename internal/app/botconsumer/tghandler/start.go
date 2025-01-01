package tghandler

import (
	"context"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/internal/tgbot"
)

func (t *TGHandler) Start(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	sender.SendMessage(ctx, msg.ChatID, "Type /order for requesting an order")

	return nil
}
