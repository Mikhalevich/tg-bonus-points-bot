package buttonrespository

import (
	"bytes"
	"context"
	"encoding/gob"
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

func encodeButton(b button.Button) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(b); err != nil {
		return nil, fmt.Errorf("gob encode: %w", err)
	}

	return buf.Bytes(), nil
}
