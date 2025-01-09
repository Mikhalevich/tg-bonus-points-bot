package customer

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/flag"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Customer) RefreshOrder(
	ctx context.Context,
	info msginfo.Info,
	orderID order.ID,
) error {
	assemblingOrder, err := c.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		if c.repository.IsNotFoundError(err) {
			c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, "no such order")
			return nil
		}

		return fmt.Errorf("get order by id: %w", err)
	}

	if !assemblingOrder.IsAssembling() {
		c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID,
			fmt.Sprintf("order in %s state", assemblingOrder.Status.HumanReadable()))
		return nil
	}

	categories, err := c.repository.GetCategoryProducts(ctx, product.Filter{
		Products: flag.Enabled,
		Category: flag.Enabled,
	})
	if err != nil {
		return fmt.Errorf("get products: %w", err)
	}

	buttons, err := c.makeOrderButtons(ctx, info.ChatID, assemblingOrder.ID, categories)
	if err != nil {
		return fmt.Errorf("make order buttons: %w", err)
	}

	c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, "Choose category", buttons...)

	return nil
}
