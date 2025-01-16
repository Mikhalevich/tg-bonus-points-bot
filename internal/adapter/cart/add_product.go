package cart

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Cart) AddProduct(ctx context.Context, cartID cart.ID, productID product.ProductID) error {
	exists, err := c.addProductToExistingList(ctx, makeCartProductsKey(cartID.String()), productID)
	if err != nil {
		return fmt.Errorf("add product to existing list: %w", err)
	}

	if !exists {
		return redis.Nil
	}

	return nil
}

// addProductToExistingList returns false is the list is not exists and true otherwise.
func (c *Cart) addProductToExistingList(ctx context.Context, key string, id product.ProductID) (bool, error) {
	newLen, err := c.client.RPushX(ctx, key, id.String()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, fmt.Errorf("rpushx: %w", err)
	}

	if newLen == 0 {
		return false, nil
	}

	return true, nil
}

//nolint:unused
func (c *Cart) addProductToNotExistingList(ctx context.Context, key string, id product.ProductID) error {
	if _, err := c.client.Pipelined(ctx, func(pipline redis.Pipeliner) error {
		if err := pipline.RPush(ctx, key, id.String()).Err(); err != nil {
			return fmt.Errorf("rpush: %w", err)
		}

		if err := pipline.Expire(ctx, key, c.ttl).Err(); err != nil {
			return fmt.Errorf("expire: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("piplined: %w", err)
	}

	return nil
}
