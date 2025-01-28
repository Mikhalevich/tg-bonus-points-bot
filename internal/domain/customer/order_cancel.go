package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (c *Customer) OrderCancel(ctx context.Context, chatID msginfo.ChatID, orderID order.ID) error {
	ord, err := c.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		if c.repository.IsNotFoundError(err) {
			c.sender.SendText(ctx, chatID, message.OrderNotExists())
			return nil
		}

		return fmt.Errorf("get order by id: %w", err)
	}

	if !ord.CanCancel() {
		c.sender.SendText(ctx, chatID, message.OrderStatus(ord.Status))
		return nil
	}

	canceledOrder, err := c.repository.UpdateOrderStatusByChatAndID(ctx, orderID, chatID, time.Now(),
		order.StatusCanceled, order.StatusWaitingPayment, order.StatusConfirmed)
	if err != nil {
		if c.repository.IsNotUpdatedError(err) {
			c.sender.SendText(ctx, chatID, message.OrderWithStatusNotExists(ord.Status))
			return nil
		}

		return fmt.Errorf("update order status: %w", err)
	}

	c.sender.SendText(ctx, chatID, message.OrderStatusChanged(canceledOrder.Status))

	return nil
}
