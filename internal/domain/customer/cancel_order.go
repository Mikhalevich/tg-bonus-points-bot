package customer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (c *Customer) CancelOrder(ctx context.Context, info msginfo.Info, orderID order.ID) error {
	assemblingOrder, err := c.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		if c.repository.IsNotFoundError(err) {
			c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, message.OrderNotExists())
			return nil
		}

		return fmt.Errorf("get order by id: %w", err)
	}

	if !assemblingOrder.IsSameChat(info.ChatID) {
		return errors.New("chat order is different")
	}

	if !assemblingOrder.CanCancel() {
		c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, message.OrderStatus(assemblingOrder.Status))
		return nil
	}

	canceledOrder, err := c.repository.UpdateOrderStatus(ctx, orderID, time.Now(), order.StatusCanceled,
		order.StatusAssembling, order.StatusConfirmed)
	if err != nil {
		if c.repository.IsNotUpdatedError(err) {
			c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID,
				message.OrderWithStatusNotExists(assemblingOrder.Status))
			return nil
		}

		return fmt.Errorf("update order status: %w", err)
	}

	c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, message.OrderStatusChanged(canceledOrder.Status))

	return nil
}
