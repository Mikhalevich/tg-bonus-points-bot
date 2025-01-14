package buttonrespository

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
)

func (r *ButtonRepository) SetButton(ctx context.Context, b button.Button) (button.InlineKeyboardButton, error) {
	encodedButton, err := encodeButton(b)
	if err != nil {
		return button.InlineKeyboardButton{}, fmt.Errorf("encode button: %w", err)
	}

	id := generateID()

	if err := r.client.Set(ctx, id, encodedButton, r.ttl).Err(); err != nil {
		return button.InlineKeyboardButton{}, fmt.Errorf("redis set: %w", err)
	}

	return button.InlineKeyboardButton{
		ID:      button.IDFromString(id),
		Caption: b.Caption,
	}, nil
}
