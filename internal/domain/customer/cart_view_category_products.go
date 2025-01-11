package customer

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Customer) CartViewCategoryProducts(
	ctx context.Context,
	info msginfo.Info,
	categoryID product.ID,
) error {
	categoryProducts, err := c.repository.GetProductsByCategoryID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("get products by category id: %w", err)
	}

	cartProducts, err := c.cart.GetProducts(ctx, info.ChatID)
	if err != nil {
		return fmt.Errorf("get cart products: %w", err)
	}

	buttons, err := c.makeCartProductsButtons(ctx, info.ChatID, categoryID, categoryProducts, cartProducts)
	if err != nil {
		return fmt.Errorf("make products buttons: %w", err)
	}

	c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, message.OrderProductPage(), buttons...)

	return nil
}

func (c *Customer) makeCartProductsButtons(
	ctx context.Context,
	chatID msginfo.ChatID,
	categoryID product.ID,
	categoryProducts []product.Product,
	cartProducts []port.CartItem,
) ([]button.InlineKeyboardButtonRow, error) {
	buttons := make([]button.InlineKeyboardButtonRow, 0, len(categoryProducts)+1)

	for _, v := range categoryProducts {
		title := makeProductButtonTitle(v, cartProducts)

		b, err := c.makeInlineKeyboardButton(ctx, button.AddProduct(chatID, v.ID, categoryID), title)
		if err != nil {
			return nil, fmt.Errorf("category order button: %w", err)
		}

		buttons = append(buttons, button.Row(b))
	}

	backBtn, err := c.makeInlineKeyboardButton(ctx, button.ViewCategories(chatID), message.Done())
	if err != nil {
		return nil, fmt.Errorf("back from products button: %w", err)
	}

	buttons = append(buttons, button.Row(backBtn))

	return buttons, nil
}

func makeProductButtonTitle(p product.Product, cartProducts []port.CartItem) string {
	for _, v := range cartProducts {
		if v.ProductID == p.ID {
			return fmt.Sprintf("%s [%d]", p.Title, v.Count)
		}
	}

	return p.Title
}
