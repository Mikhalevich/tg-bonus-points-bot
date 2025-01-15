package cart

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Cart) GetProducts(ctx context.Context, id cart.ID) ([]port.CartItem, error) {
	items, err := c.client.LRange(ctx, makeCartProductsKey(id.String()), 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("lrange: %w", err)
	}

	cartItems, err := convertToCartItems(combineDuplicateItems(items))
	if err != nil {
		return nil, fmt.Errorf("convert to cart items: %w", err)
	}

	return cartItems, nil
}

func combineDuplicateItems(items []string) map[string]int {
	itemsMap := make(map[string]int, len(items))

	for _, v := range items {
		if _, ok := itemsMap[v]; !ok {
			itemsMap[v] = 1
			continue
		}

		itemsMap[v]++
	}

	return itemsMap
}

func convertToCartItems(itemsMap map[string]int) ([]port.CartItem, error) {
	cartItems := make([]port.CartItem, 0, len(itemsMap))

	for id, count := range itemsMap {
		if isEmptyListPlaceholder(id) {
			continue
		}

		productID, err := product.IDFromString(id)
		if err != nil {
			return nil, fmt.Errorf("make id from string: %w", err)
		}

		cartItems = append(cartItems, port.CartItem{
			ProductID: productID,
			Count:     count,
		})
	}

	return cartItems, nil
}
