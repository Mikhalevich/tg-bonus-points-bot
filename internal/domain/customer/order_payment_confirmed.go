package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (c *Customer) OrderPaymentConfirmed(ctx context.Context, chatID msginfo.ChatID, orderID order.ID) error {
	if _, err := c.repository.UpdateOrderStatusByChatAndID(
		ctx,
		orderID,
		chatID,
		time.Now(),
		order.StatusConfirmed,
		order.StatusPaymentInProgress,
	); err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	return nil
}
