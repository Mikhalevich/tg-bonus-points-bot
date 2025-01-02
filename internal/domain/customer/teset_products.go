package customer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

func (c *Customer) TestProducts(ctx context.Context, chatID msginfo.ChatID) error {
	cat, err := c.repository.GetCategoryProducts(ctx)
	if err != nil {
		return fmt.Errorf("get category products: %w", err)
	}

	b, err := json.Marshal(cat)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	c.sender.SendText(ctx, chatID, string(b))
	return nil
}
