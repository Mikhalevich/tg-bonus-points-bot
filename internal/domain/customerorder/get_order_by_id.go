package customerorder

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (c *CustomerOrder) GetOrderByID(ctx context.Context, chatID msginfo.ChatID, orderID order.ID) error {
	ord, err := c.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		if c.repository.IsNotFoundError(err) {
			c.sender.SendText(ctx, chatID, message.InvalidOrder())
			return nil
		}

		return fmt.Errorf("get order by chat_id: %w", err)
	}

	if !ord.IsSameChat(chatID) {
		c.sender.SendText(ctx, chatID, message.InvalidOrder())
		return nil
	}

	productsInfo, err := c.repository.GetProductsByIDs(ctx, ord.ProductIDs(), ord.CurrencyID)
	if err != nil {
		return fmt.Errorf("get products by ids: %w", err)
	}

	c.sender.SendText(ctx, chatID, formatOrder(ord, productsInfo, 0, c.sender.EscapeMarkdown))

	return nil
}
