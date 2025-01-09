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
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	orderID order.ID,
) error {
	assemblingOrder, err := c.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("get order by id: %w", err)
	}

	if assemblingOrder.Status != order.StatusAssembling {
		c.sender.EditTextMessage(ctx, chatID, messageID,
			fmt.Sprintf("order has different state: %s", assemblingOrder.Status.HumanReadable()))
		return nil
	}

	categories, err := c.repository.GetCategoryProducts(ctx, product.Filter{
		Products: flag.Enabled,
		Category: flag.Enabled,
	})
	if err != nil {
		return fmt.Errorf("get products: %w", err)
	}

	buttons, err := c.makeOrderButtons(ctx, chatID, assemblingOrder.ID, categories)
	if err != nil {
		return fmt.Errorf("make order buttons: %w", err)
	}

	c.sender.EditTextMessage(ctx, chatID, messageID, "Choose category", buttons...)

	return nil
}
