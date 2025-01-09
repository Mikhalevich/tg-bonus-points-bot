package customer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (c *Customer) ConfirmOrder(ctx context.Context, info msginfo.Info, orderID order.ID) error {
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

	if !assemblingOrder.IsAssembling() {
		c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, message.OrderStatus(assemblingOrder.Status))
		return nil
	}

	confirmedOrder, err := c.repository.UpdateOrderStatus(ctx, orderID, time.Now(),
		order.StatusConfirmed, order.StatusAssembling)
	if err != nil {
		if c.repository.IsNotUpdatedError(err) {
			c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID,
				message.OrderWithStatusNotExists(assemblingOrder.Status))
			return nil
		}

		return fmt.Errorf("update order status: %w", err)
	}

	cancelBtn, err := c.makeInlineKeyboardButton(ctx, button.CancelOrder(confirmedOrder.ChatID, orderID), message.Cancel())
	if err != nil {
		return fmt.Errorf("cancel order button: %w", err)
	}

	png, err := c.qrCode.GeneratePNG(orderID.String())
	if err != nil {
		return fmt.Errorf("qrcode generate png: %w", err)
	}

	if err := c.sender.SendPNGMarkdown(
		ctx,
		confirmedOrder.ChatID,
		formatOrder(confirmedOrder, c.sender.EscapeMarkdown),
		png,
		button.Row(cancelBtn),
	); err != nil {
		return fmt.Errorf("send png: %w", err)
	}

	c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, "completed")

	return nil
}
