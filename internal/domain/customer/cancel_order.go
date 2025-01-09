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

func (c *Customer) CancelOrderSendMessage(ctx context.Context, chatID msginfo.ChatID, orderID order.ID) error {
	return c.cancelOrder(
		ctx,
		msginfo.Info{
			ChatID: chatID,
		},
		orderID)
}

func (c *Customer) CancelOrderEditMessage(ctx context.Context, info msginfo.Info, orderID order.ID) error {
	return c.cancelOrder(ctx, info, orderID)
}

func (c *Customer) cancelOrder(ctx context.Context, info msginfo.Info, orderID order.ID) error {
	assemblingOrder, err := c.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		if c.repository.IsNotFoundError(err) {
			c.sendOrEditMessage(ctx, info, message.OrderNotExists())
			return nil
		}

		return fmt.Errorf("get order by id: %w", err)
	}

	if !assemblingOrder.IsSameChat(info.ChatID) {
		return errors.New("chat order is different")
	}

	if !assemblingOrder.CanCancel() {
		c.sendOrEditMessage(ctx, info, message.OrderStatus(assemblingOrder.Status))
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

	c.sendOrEditMessage(ctx, info, message.OrderStatusChanged(canceledOrder.Status))

	return nil
}

func (c *Customer) sendOrEditMessage(ctx context.Context, info msginfo.Info, msg string) {
	if info.MessageID.Int() != 0 {
		c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, msg)
		return
	}

	c.sender.SendText(ctx, info.ChatID, msg)
}
