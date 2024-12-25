package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (c *Customer) TestCompleteOrder(ctx context.Context, info msginfo.Info) error {
	activeOrder, err := c.repository.GetOrderByChatIDAndStatus(
		ctx,
		info.ChatID,
		order.StatusCreated,
		order.StatusInProgress,
		order.StatusReady,
	)

	if err != nil {
		if c.repository.IsNotFoundError(err) {
			c.sender.ReplyText(ctx, info.ChatID, info.MessageID, "no active orders")
			return nil
		}

		return fmt.Errorf("get order by chat_id: %w", err)
	}

	if err := c.repository.UpdateOrderStatus(
		ctx,
		activeOrder.ID,
		time.Now(),
		order.StatusCompleted,
		order.StatusCreated,
		order.StatusInProgress,
		order.StatusReady,
	); err != nil {
		if c.repository.IsNotUpdatedError(err) {
			c.sender.ReplyText(ctx, info.ChatID, info.MessageID, "order already completed")
			return nil
		}

		return fmt.Errorf("update order status: %w", err)
	}

	c.sender.ReplyText(ctx, info.ChatID, info.MessageID,
		fmt.Sprintf("order %s completed successfully", activeOrder.ID))

	return nil
}
