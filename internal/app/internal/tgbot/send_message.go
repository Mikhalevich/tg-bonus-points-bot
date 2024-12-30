package tgbot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
)

func (t *TGBot) SendMessage(ctx context.Context, chatID int64, msg string) error {
	if _, err := t.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   msg,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
