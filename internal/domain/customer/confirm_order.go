package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (c *Customer) ConfirmOrder(ctx context.Context, chatID msginfo.ChatID, orderID order.ID) error {
	assemblingOrder, err := c.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		if c.repository.IsNotFoundError(err) {
			c.sender.SendText(ctx, chatID, "Order not found")
			return nil
		}

		return fmt.Errorf("get order by id: %w", err)
	}

	if !assemblingOrder.IsSameChat(chatID) {
		c.sender.SendText(ctx, chatID, "Order permission failure")
		return nil
	}

	if !assemblingOrder.CanConfirm() {
		c.sender.SendTextMarkdown(ctx, chatID,
			fmt.Sprintf("order cannot be confirmed from *%s* state", assemblingOrder.Status.HumanReadable()))
		return nil
	}

	confirmedOrder, err := c.repository.UpdateOrderStatus(ctx, orderID, time.Now(),
		order.StatusConfirmed, order.StatusAssembling)
	if err != nil {
		if c.repository.IsNotUpdatedError(err) {
			c.sender.SendText(ctx, chatID, "Order cannot be confirmed")
			return nil
		}

		return fmt.Errorf("update order status: %w", err)
	}

	cancelBtn, err := c.makeInlineKeyboardButton(ctx, button.CancelOrder(confirmedOrder.ChatID, orderID), "Cancel")
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
		cancelBtn,
	); err != nil {
		return fmt.Errorf("send png: %w", err)
	}

	return nil
}
