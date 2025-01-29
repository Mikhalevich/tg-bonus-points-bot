package customer

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (c *Customer) CartViewCategories(
	ctx context.Context,
	info msginfo.Info,
	cartID cart.ID,
) error {
	cartItems, err := c.cart.GetProducts(ctx, cartID)
	if err != nil {
		if c.cart.IsNotFoundError(err) {
			c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, message.CartOrderUnavailable())
			return nil
		}

		return fmt.Errorf("cart items: %w", err)
	}

	orderedProducts, err := c.orderedProducts(ctx, cartItems)
	if err != nil {
		return fmt.Errorf("cart products: %w", err)
	}

	categories, err := c.repository.GetCategories(ctx)
	if err != nil {
		return fmt.Errorf("get products: %w", err)
	}

	buttons, err := c.makeCartCategoriesButtons(
		ctx,
		info.ChatID,
		cartID,
		categories,
		orderedProducts,
		currencyFromProducts(orderedProducts),
	)
	if err != nil {
		return fmt.Errorf("make order buttons: %w", err)
	}

	c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, message.OrderCategoryPage(), buttons...)

	return nil
}

func currencyFromProducts(products []order.OrderedProduct) currency.Currency {
	if len(products) == 0 {
		return currency.Currency{}
	}

	return products[0].Product.Currency
}
