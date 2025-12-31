package messageprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
)

func (m *MessageProcessor) SetButton(
	ctx context.Context,
	btn button.Button,
) (button.InlineKeyboardButton, error) {
	if btn.Pay {
		return button.InlineKeyboardButton{
			Caption: btn.Caption,
			Pay:     true,
		}, nil
	}

	if err := m.buttonRepository.SetButton(ctx, btn); err != nil {
		return button.InlineKeyboardButton{}, fmt.Errorf("button repository set button: %w", err)
	}

	return button.InlineKeyboardButton{
		ID:      btn.ID,
		Caption: btn.Caption,
	}, nil
}
