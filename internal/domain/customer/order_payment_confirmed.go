package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (c *Customer) OrderPaymentConfirmed(
	ctx context.Context,
	chatID msginfo.ChatID,
	orderID order.ID,
	currency string,
	totalAmount int,
) error {
	ord, err := c.repository.UpdateOrderStatusByChatAndID(
		ctx,
		orderID,
		chatID,
		time.Now(),
		order.StatusConfirmed,
		order.StatusPaymentInProgress,
	)
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	if err := c.sendOrderQRImage(ctx, chatID, ord); err != nil {
		return fmt.Errorf("send order qr: %w", err)
	}

	return nil
}

func (c *Customer) sendOrderQRImage(ctx context.Context, chatID msginfo.ChatID, ord *order.Order) error {
	png, err := c.qrCode.GeneratePNG(ord.ID.String())
	if err != nil {
		return fmt.Errorf("qrcode generate png: %w", err)
	}

	productsInfo, err := c.repository.GetProductsByIDs(ctx, ord.ProductIDs(), ord.CurrencyID)
	if err != nil {
		return fmt.Errorf("get products by ids: %w", err)
	}

	if err := c.sender.SendPNGMarkdown(
		ctx,
		chatID,
		formatOrder(ord, productsInfo, c.sender.EscapeMarkdown),
		png,
	); err != nil {
		return fmt.Errorf("send png: %w", err)
	}

	return nil
}
