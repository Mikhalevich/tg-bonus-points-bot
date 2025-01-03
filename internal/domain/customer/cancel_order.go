package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (c *Customer) CancelOrder(ctx context.Context, id order.ID) error {
	activeOrder, err := c.repository.GetOrderByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get order by id: %w", err)
	}

	if !isOrderCancelable(activeOrder.Status) {
		c.sender.SendTextMarkdown(ctx, activeOrder.ChatID,
			fmt.Sprintf("order cannot be canceled in *%s* state", activeOrder.Status.HumanReadable()))
		return nil
	}

	canceledOrder, err := c.repository.UpdateOrderStatus(ctx, id, time.Now(), order.StatusCanceled, order.StatusCreated)
	if err != nil {
		if c.repository.IsNotUpdatedError(err) {
			c.sender.SendText(ctx, activeOrder.ChatID, "order cannot be canceled")
			return nil
		}

		return fmt.Errorf("update order status: %w", err)
	}

	c.sender.SendTextMarkdown(ctx, canceledOrder.ChatID,
		fmt.Sprintf("order canceled successfully\n%s", formatOrder(canceledOrder, c.sender.EscapeMarkdown)))

	return nil
}

func isOrderCancelable(s order.Status) bool {
	return s == order.StatusCreated
}
