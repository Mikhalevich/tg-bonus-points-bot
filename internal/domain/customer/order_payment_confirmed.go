package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
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
	now := time.Now()

	position, err := c.dailyPosition.Position(ctx, now)
	if err != nil {
		return fmt.Errorf("daily position: %w", err)
	}

	ord, err := c.repository.UpdateOrderByChatAndID(
		ctx,
		orderID,
		chatID,
		port.UpdateOrderData{
			Status:              order.StatusConfirmed,
			StatusOperationTime: now,
			VerificationCode:    c.codeGenerator.Generate(),
			DailyPosition:       position,
		},
		order.StatusPaymentInProgress,
	)
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	queuePosition := c.orderQueuePosition(ctx, ord)

	if err := c.sendOrderQRImage(ctx, chatID, ord, queuePosition); err != nil {
		return fmt.Errorf("send order qr: %w", err)
	}

	return nil
}

func (c *Customer) sendOrderQRImage(
	ctx context.Context,
	chatID msginfo.ChatID,
	ord *order.Order,
	queuePosition int,
) error {
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
		formatOrder(ord, productsInfo, queuePosition, c.sender.EscapeMarkdown),
		png,
	); err != nil {
		return fmt.Errorf("send png: %w", err)
	}

	return nil
}
