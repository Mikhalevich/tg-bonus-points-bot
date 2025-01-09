package customer

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Customer) ViewCategoryProducts(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	categoryID product.ID,
) error {
	products, err := c.repository.GetProductsByCategoryID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("get products by category id: %w", err)
	}

	buttons, err := c.makeProductsButtons(ctx, chatID, products)
	if err != nil {
		return fmt.Errorf("make products buttons: %w", err)
	}

	c.sender.EditTextMessage(ctx, chatID, messageID, "choose product", buttons...)

	return nil
}

func (c *Customer) makeProductsButtons(
	ctx context.Context,
	chatID msginfo.ChatID,
	products []product.Product,
) ([]button.InlineKeyboardButtonRow, error) {
	buttons := make([]button.InlineKeyboardButtonRow, 0, len(products)+1)

	for _, v := range products {
		b, err := c.makeInlineKeyboardButton(ctx, button.Product(chatID, v.ID), v.Title)
		if err != nil {
			return nil, fmt.Errorf("category order button: %w", err)
		}

		buttons = append(buttons, button.Row(b))
	}

	backBtn, err := c.makeInlineKeyboardButton(ctx, button.BackToOrder(chatID, 0), "Back")
	if err != nil {
		return nil, fmt.Errorf("back from products button: %w", err)
	}

	buttons = append(buttons, button.Row(backBtn))

	return buttons, nil
}
