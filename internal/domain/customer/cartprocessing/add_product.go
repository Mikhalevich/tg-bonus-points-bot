package cartprocessing

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func (c *CartProcessing) AddProduct(
	ctx context.Context,
	info msginfo.Info,
	cartID cart.ID,
	categoryID product.CategoryID,
	productID product.ProductID,
	currencyID currency.ID,
) error {
	if err := c.cart.AddProduct(ctx, cartID, cart.CartProduct{
		ProductID:  productID,
		CategoryID: categoryID,
		Count:      1,
	}); err != nil {
		if c.cart.IsNotFoundError(err) {
			c.editPlainText(ctx, info.ChatID, info.MessageID, message.CartOrderUnavailable())

			return nil
		}

		return fmt.Errorf("add product to cart: %w", err)
	}

	if err := c.ViewCategoryProducts(ctx, info, cartID, categoryID, currencyID); err != nil {
		return fmt.Errorf("view category products: %w", err)
	}

	return nil
}
