package customer

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Customer) CartAddProduct(
	ctx context.Context,
	info msginfo.Info,
	cartID cart.ID,
	categoryID product.CategoryID,
	productID product.ProductID,
) error {
	if err := c.cart.AddProduct(ctx, cartID, productID); err != nil {
		if c.cart.IsNotFoundError(err) {
			c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, message.CartOrderUnavailable())
			return nil
		}

		return fmt.Errorf("add product to cart: %w", err)
	}

	if err := c.CartViewCategoryProducts(ctx, info, cartID, categoryID); err != nil {
		return fmt.Errorf("view category products: %w", err)
	}

	return nil
}
