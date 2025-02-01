package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (c *Customer) OrderCancel(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	orderID order.ID,
	isTextMsg bool,
) error {
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

	if _, err := c.repository.UpdateOrderStatusByChatAndID(ctx, orderID, chatID, time.Now(),
		order.StatusCanceled, order.StatusWaitingPayment, order.StatusConfirmed); err != nil {
		if c.repository.IsNotUpdatedError(err) {
			c.sender.SendText(ctx, chatID, message.OrderWithStatusNotExists(ord.Status))
			return nil
		}

		return fmt.Errorf("update order status: %w", err)
	}

	c.editOridinOrderMessage(ctx, chatID, messageID, isTextMsg)

	return nil
}

func (c *Customer) editOridinOrderMessage(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	isTextMsg bool,
) {
	if isTextMsg {
		c.sender.EditTextMessage(ctx, chatID, messageID, message.OrderCanceled())
		return
	}

	c.sender.DeleteMessage(ctx, chatID, messageID)
	c.sender.SendText(ctx, chatID, message.OrderCanceled())
}
