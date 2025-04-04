package buttonrespository

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
)

func (r *ButtonRepository) SetButton(ctx context.Context, btn button.Button) (button.InlineKeyboardButton, error) {
	encodedButton, err := encodeButton(btn)
	if err != nil {
		return button.InlineKeyboardButton{}, fmt.Errorf("encode button: %w", err)
	}

	btnID := generateID()

	if err := r.client.Set(ctx, btnID, encodedButton, r.ttl).Err(); err != nil {
		return button.InlineKeyboardButton{}, fmt.Errorf("redis set: %w", err)
	}

	return button.InlineKeyboardButton{
		ID:      button.IDFromString(btnID),
		Caption: btn.Caption,
	}, nil
}
