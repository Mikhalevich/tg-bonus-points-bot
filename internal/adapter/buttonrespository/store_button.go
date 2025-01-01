package buttonrespository

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
)

func (r *ButtonRepository) StoreButton(ctx context.Context, b *button.Button) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(b); err != nil {
		return fmt.Errorf("gob encode: %w", err)
	}

	if err := r.client.Set(ctx, b.ID.String(), buf.Bytes(), r.ttl).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	return nil
}
