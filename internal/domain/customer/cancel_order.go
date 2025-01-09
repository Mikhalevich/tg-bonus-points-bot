package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (c *Customer) CancelOrder(ctx context.Context, info msginfo.Info, orderID order.ID) error {
	activeOrder, err := c.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		if c.repository.IsNotFoundError(err) {
			c.sender.SendText(ctx, info.ChatID, "Order not found")
			return nil
		}

		return fmt.Errorf("get order by id: %w", err)
	}

	if !activeOrder.IsSameChat(info.ChatID) {
		c.sender.SendText(ctx, info.ChatID, "Order permission failure")
		return nil
	}

	if !activeOrder.CanCancel() {
		c.sender.SendTextMarkdown(ctx, info.ChatID,
			fmt.Sprintf("order cannot be canceled in *%s* state", activeOrder.Status.HumanReadable()))
		return nil
	}

	canceledOrder, err := c.repository.UpdateOrderStatus(ctx, orderID, time.Now(), order.StatusCanceled,
		order.StatusAssembling, order.StatusConfirmed)
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
