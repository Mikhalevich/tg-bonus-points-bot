package cartprocessing

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (c *CartProcessing) Cancel(ctx context.Context, info msginfo.Info, cartID cart.ID) error {
	if err := c.cart.Clear(ctx, info.ChatID, cartID); err != nil {
		if c.cart.IsNotFoundError(err) {
			c.editPlainText(ctx, info.ChatID, info.MessageID, message.CartOrderUnavailable())

			return nil
		}

		return fmt.Errorf("clear cart: %w", err)
	}

	if err := c.sender.DeleteMessage(ctx, info.ChatID, info.MessageID); err != nil {
		return fmt.Errorf("delete message: %w", err)
	}

	return nil
}
