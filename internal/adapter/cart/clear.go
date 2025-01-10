package cart

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

func (c *Cart) Clear(ctx context.Context, chatID msginfo.ChatID) error {
	if err := c.client.Del(ctx, makeKey(chatID)).Err(); err != nil {
		return fmt.Errorf("del: %w", err)
	}

	return nil
}
