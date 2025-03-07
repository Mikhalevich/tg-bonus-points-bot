package orderaction

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (o *OrderAction) Cancel(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	orderID order.ID,
	isTextMsg bool,
) error {
	ord, err := o.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		if o.repository.IsNotFoundError(err) {
			o.sender.SendText(ctx, chatID, message.OrderNotExists())
			return nil
		}

		return fmt.Errorf("get order by id: %w", err)
	}

	if !ord.CanCancel() {
		o.sender.SendText(ctx, chatID, message.OrderStatus(ord.Status))
		return nil
	}

	if _, err := o.repository.UpdateOrderStatusByChatAndID(ctx, orderID, chatID, o.timeProvider.Now(),
		order.StatusCanceled, order.StatusWaitingPayment, order.StatusConfirmed); err != nil {
		if o.repository.IsNotUpdatedError(err) {
			o.sender.SendText(ctx, chatID, message.OrderWithStatusNotExists(ord.Status))
			return nil
		}

		return fmt.Errorf("update order status: %w", err)
	}

	o.editOridinOrderMessage(ctx, chatID, messageID, isTextMsg)

	return nil
}

func (o *OrderAction) editOridinOrderMessage(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	isTextMsg bool,
) {
	if isTextMsg {
		o.sender.EditTextMessage(ctx, chatID, messageID, message.OrderCanceled())
		return
	}

	o.sender.DeleteMessage(ctx, chatID, messageID)
	o.sender.SendText(ctx, chatID, message.OrderCanceled())
}
