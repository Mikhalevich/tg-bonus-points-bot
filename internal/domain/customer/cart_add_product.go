package customer

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Customer) CartAddProduct(
	ctx context.Context,
	info msginfo.Info,
	categoryID product.ID,
	productID product.ID,
) error {
	if err := c.cart.AddProduct(ctx, info.ChatID, productID); err != nil {
		return fmt.Errorf("add product to cart: %w", err)
	}

	if err := c.CartViewCategoryProducts(ctx, info, categoryID); err != nil {
		return fmt.Errorf("view category products: %w", err)
	}

	return nil
}
