package customer

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Customer) CartViewCategoryProducts(
	ctx context.Context,
	info msginfo.Info,
	cartID cart.ID,
	categoryID product.CategoryID,
	currencyID currency.ID,
) error {
	cartProducts, err := c.cart.GetProducts(ctx, cartID)
	if err != nil {
		if c.cart.IsNotFoundError(err) {
			c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, message.CartOrderUnavailable())
			return nil
		}

		return fmt.Errorf("get cart products: %w", err)
	}

	categoryProducts, err := c.repository.GetProductsByCategoryID(ctx, categoryID, currencyID)
	if err != nil {
		return fmt.Errorf("get products by category id: %w", err)
	}

	buttons, err := c.makeCartProductsButtons(
		ctx,
		info.ChatID,
		cartID,
		categoryID,
		categoryProducts,
		cartProducts,
		currencyID,
	)
	if err != nil {
		return fmt.Errorf("make products buttons: %w", err)
	}

	c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, message.OrderProductPage(), buttons...)

	return nil
}

func (c *Customer) makeCartProductsButtons(
	ctx context.Context,
	chatID msginfo.ChatID,
	cartID cart.ID,
	categoryID product.CategoryID,
	categoryProducts []product.Product,
	cartProducts []cart.CartProduct,
	currencyID currency.ID,
) ([]button.InlineKeyboardButtonRow, error) {
	buttons := make([]button.ButtonRow, 0, len(categoryProducts)+1)

	for _, v := range categoryProducts {
		title := makeProductButtonTitle(v, cartProducts)
		btn, err := button.CartAddProduct(chatID, title, cartID, v.ID, categoryID, currencyID)

		if err != nil {
			return nil, fmt.Errorf("add product button: %w", err)
		}

		buttons = append(buttons, button.Row(btn))
	}

	viewCategoriesBtn, err := button.CartViewCategories(chatID, message.Done(), cartID, currencyID)
	if err != nil {
		return nil, fmt.Errorf("cart view categories button: %w", err)
	}

	buttons = append(buttons, button.Row(viewCategoriesBtn))

	inlineButtons, err := c.buttonRepository.SetButtonRows(ctx, buttons...)
	if err != nil {
		return nil, fmt.Errorf("set button rows: %w", err)
	}

	return inlineButtons, nil
}

func makeProductButtonTitle(p product.Product, cartProducts []cart.CartProduct) string {
	for _, v := range cartProducts {
		if v.ProductID == p.ID {
			return fmt.Sprintf("%s %s [x%d %s]",
				p.Title,
				p.Currency.FormatPrice(p.Price),
				v.Count,
				p.Currency.FormatPrice(p.Price*v.Count),
			)
		}
	}

	return fmt.Sprintf("%s %s",
		p.Title,
		p.Currency.FormatPrice(p.Price),
	)
}
