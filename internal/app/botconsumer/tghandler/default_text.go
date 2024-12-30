package tghandler

import (
	"context"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/internal/tgbot"
)

func (t *TGHandler) DefaultText(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	// skip message.
	return nil
}
