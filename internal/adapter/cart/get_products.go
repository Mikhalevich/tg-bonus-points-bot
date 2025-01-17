package cart

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/cart"
)

func (c *Cart) GetProducts(ctx context.Context, id cart.ID) ([]cart.CartProduct, error) {
	items, err := c.client.LRange(ctx, makeCartProductsKey(id.String()), 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("lrange: %w", err)
	}

	if len(items) == 0 {
		return nil, redis.Nil
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

func convertToCartItems(itemsMap map[string]int) ([]cart.CartProduct, error) {
	cartItems := make([]cart.CartProduct, 0, len(itemsMap))

	for encodedProduct, count := range itemsMap {
		if isEmptyListPlaceholder(encodedProduct) {
			continue
		}

		p, err := decodeCartProduct(encodedProduct)
		if err != nil {
			return nil, fmt.Errorf("decode cart product: %w", err)
		}

		cartItems = append(cartItems, cart.CartProduct{
			ProductID:  p.ProductID,
			CategoryID: p.CategoryID,
			Count:      count,
		})
	}

	return cartItems, nil
}
