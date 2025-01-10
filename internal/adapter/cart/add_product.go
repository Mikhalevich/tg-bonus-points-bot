package cart

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Cart) AddProduct(ctx context.Context, chatID msginfo.ChatID, productID product.ID) error {
	key := makeKey(chatID)

	exists, err := c.addProductToExistingList(ctx, key, productID)
	if err != nil {
		return fmt.Errorf("add product to existing list: %w", err)
	}

	if exists {
		return nil
	}

	if err := c.addProductToNotExistingList(ctx, key, productID); err != nil {
		return fmt.Errorf("add product to not existing list: %w", err)
	}

	return nil
}

func makeKey(chatID msginfo.ChatID) string {
	return fmt.Sprintf("cart:%d", chatID.Int64())
}

// addProductToExistingList returns false is the list is not exists and true otherwise.
func (c *Cart) addProductToExistingList(ctx context.Context, key string, id product.ID) (bool, error) {
	if err := c.client.RPushX(ctx, key, id.String()).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, fmt.Errorf("rpushx: %w", err)
	}

	return true, nil
}

func (c *Cart) addProductToNotExistingList(ctx context.Context, key string, id product.ID) error {
	if _, err := c.client.Pipelined(ctx, func(pipline redis.Pipeliner) error {
		if err := pipline.RPush(ctx, key, id).Err(); err != nil {
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
