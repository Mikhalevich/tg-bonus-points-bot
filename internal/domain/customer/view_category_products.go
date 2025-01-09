package customer

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Customer) ViewCategoryProducts(ctx context.Context, chatID msginfo.ChatID, categoryID product.ID) error {
	_, err := c.repository.GetProductsByCategoryID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("get products by category id: %w", err)
	}

	return nil
}
