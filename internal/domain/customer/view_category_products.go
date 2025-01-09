package customer

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Customer) ViewCategoryProducts(
	ctx context.Context,
	info msginfo.Info,
	orderID order.ID,
	categoryID product.ID,
) error {
	assemblingOrder, err := c.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("get order by id: %w", err)
	}

	if assemblingOrder.Status != order.StatusAssembling {
		return fmt.Errorf("invalid order status: %s", assemblingOrder.Status.HumanReadable())
	}

	products, err := c.repository.GetProductsByCategoryID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("get products by category id: %w", err)
	}

	buttons, err := c.makeProductsButtons(ctx, info.ChatID, orderID, products)
	if err != nil {
		return fmt.Errorf("make products buttons: %w", err)
	}

	c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, "choose product", buttons...)

	return nil
}

func (c *Customer) makeProductsButtons(
	ctx context.Context,
	chatID msginfo.ChatID,
	orderID order.ID,
	products []product.Product,
) ([]button.InlineKeyboardButtonRow, error) {
	buttons := make([]button.InlineKeyboardButtonRow, 0, len(products)+1)

	for _, v := range products {
		b, err := c.makeInlineKeyboardButton(ctx, button.Product(chatID, orderID, v.ID), v.Title)
		if err != nil {
			return nil, fmt.Errorf("category order button: %w", err)
		}

		buttons = append(buttons, button.Row(b))
	}

	backBtn, err := c.makeInlineKeyboardButton(ctx, button.BackToOrder(chatID, orderID), "Back")
	if err != nil {
		return nil, fmt.Errorf("back from products button: %w", err)
	}

	buttons = append(buttons, button.Row(backBtn))

	return buttons, nil
}
