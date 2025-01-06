package customer

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
)

func (c *Customer) makeInlineKeyboardButton(
	ctx context.Context,
	btn button.Button,
	caption string,
) (button.InlineKeyboardButton, error) {
	if err := c.buttonRepository.StoreButton(ctx, &btn); err != nil {
		return button.InlineKeyboardButton{}, fmt.Errorf("store button: %w", err)
	}

	return button.InlineKeyboardButton{
		ID:      btn.ID,
		Caption: caption,
	}, nil
}
