package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (c *Customer) OrderPaymentConfirmed(
	ctx context.Context,
	orderID order.ID,
	currency string,
	totalAmount int) error {
	if _, err := c.repository.UpdateOrderStatus(
		ctx,
		orderID,
		time.Now(),
		order.StatusConfirmed,
		order.StatusPaymentInProgress,
	); err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	return nil
}
