package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/internal/tgbot"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (t *TGHandler) DefaultHandler(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if msg.Payment.IsCheckout {
		if err := t.processCheckoutPayment(ctx, msg.Payment); err != nil {
			return fmt.Errorf("process payment: %w", err)
		}
	}

	if msg.Payment.IsSuccessful {
		if err := t.processSuccessfulPayment(ctx, msg.Payment); err != nil {
			return fmt.Errorf("process payment: %w", err)
		}
	}

	return nil
}

func (t *TGHandler) processCheckoutPayment(ctx context.Context, payment tgbot.Payment) error {
	orderID, err := order.IDFromString(payment.InvoicePayload)
	if err != nil {
		return fmt.Errorf("invalid order id: %w", err)
	}

	if err := t.orderProcessor.OrderSetPaymentInProgress(
		ctx,
		payment.ID,
		orderID,
		payment.Currency,
		payment.TotalAmount,
	); err != nil {
		return fmt.Errorf("set payment in progress: %w", err)
	}

	return nil
}

func (t *TGHandler) processSuccessfulPayment(ctx context.Context, payment tgbot.Payment) error {
	orderID, err := order.IDFromString(payment.InvoicePayload)
	if err != nil {
		return fmt.Errorf("invalid order id: %w", err)
	}

	if err := t.orderProcessor.OrderPaymentConfirmed(ctx, orderID, payment.Currency, payment.TotalAmount); err != nil {
		return fmt.Errorf("payment confirmed: %w", err)
	}

	return nil
}
