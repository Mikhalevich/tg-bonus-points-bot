package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (c *Customer) OrderSetPaymentInProgress(
	ctx context.Context,
	paymentID string,
	orderID order.ID,
	currency string,
	totalAmount int,
) error {
	ord, err := c.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("get order by id: %w", err)
	}

	if ord.Status != order.StatusWaitingPayment {
		return fmt.Errorf("invalid order status expected: %s actual: %s", order.StatusWaitingPayment, ord.Status)
	}

	if ord.TotalPrice() != totalAmount {
		return fmt.Errorf("invalid total amount")
	}

	if _, err := c.repository.UpdateOrderStatus(
		ctx,
		orderID,
		time.Now(),
		order.StatusPaymentInProgress,
		order.StatusWaitingPayment,
	); err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	return nil
}
