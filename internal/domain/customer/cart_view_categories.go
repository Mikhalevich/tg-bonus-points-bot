package customer

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

func (c *Customer) CartViewCategories(
	ctx context.Context,
	info msginfo.Info,
	cartID cart.ID,
) error {
	categories, err := c.repository.GetCategories(ctx)
	if err != nil {
		return fmt.Errorf("get products: %w", err)
	}

	buttons, err := c.makeCartCategoriesButtons(ctx, info.ChatID, cartID, categories)
	if err != nil {
		return fmt.Errorf("make order buttons: %w", err)
	}

	c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, message.OrderCategoryPage(), buttons...)

	return nil
}
