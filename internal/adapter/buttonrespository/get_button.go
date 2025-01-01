package buttonrespository

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
)

func (r *ButtonRepository) GetButton(ctx context.Context, id button.ID) (*button.Button, error) {
	b, err := r.client.GetDel(ctx, id.String()).Bytes()
	if err != nil {
		return nil, fmt.Errorf("redis get: %w", err)
	}

	var btn button.Button
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&btn); err != nil {
		return nil, fmt.Errorf("gob decode: %w", err)
	}

	return &btn, nil
}
