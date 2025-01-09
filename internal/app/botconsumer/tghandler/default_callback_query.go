package tghandler

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/internal/tgbot"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/perror"
)

//nolint:cyclop
func (t *TGHandler) DefaultCallbackQuery(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	btn, err := t.orderProcessor.GetButton(ctx, button.IDFromString(msg.Data))
	if err != nil {
		if perror.IsType(err, perror.TypeNotFound) {
			sender.SendMessage(ctx, msg.ChatID, "Action already executed or expired")
			return nil
		}

		return fmt.Errorf("get button: %w", err)
	}

	if btn.ChatID.Int64() != msg.ChatID {
		return fmt.Errorf("chat not match button: %d msg: %d", btn.ChatID.Int64(), msg.ChatID)
	}

	switch btn.Operation {
	case button.OperationCancelOrder:
		if err := t.cancelOrder(ctx, btn); err != nil {
			return fmt.Errorf("cancel order: %w", err)
		}

	case button.OperationConfirmOrder:
		if err := t.confirmOrder(ctx, btn); err != nil {
			return fmt.Errorf("confirm order: %w", err)
		}

	case button.OperationViewCategory:
		if err := t.ViewProducts(ctx, msginfo.MessageIDFromInt(msg.MessageID), btn); err != nil {
			return fmt.Errorf("view products: %w", err)
		}

	case button.OperationProduct:
		return errors.New("not implemented")

	case button.OperationBackToOrder:
		return errors.New("not implemented")
	}

	return nil
}

func (t *TGHandler) cancelOrder(ctx context.Context, btn *button.Button) error {
	orderID, err := btn.OrderID()
	if err != nil {
		return fmt.Errorf("invalid order id: %w", err)
	}

	if err := t.orderProcessor.CancelOrder(ctx, btn.ChatID, orderID); err != nil {
		return fmt.Errorf("cancel order: %w", err)
	}

	return nil
}

func (t *TGHandler) confirmOrder(ctx context.Context, btn *button.Button) error {
	orderID, err := btn.OrderID()
	if err != nil {
		return fmt.Errorf("invalid order id: %w", err)
	}

	if err := t.orderProcessor.ConfirmOrder(ctx, btn.ChatID, orderID); err != nil {
		return fmt.Errorf("confirm order: %w", err)
	}

	return nil
}

func (t *TGHandler) ViewProducts(ctx context.Context, messageID msginfo.MessageID, btn *button.Button) error {
	payload, err := btn.ViewCategoryPayload()
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.orderProcessor.ViewCategoryProducts(
		ctx,
		btn.ChatID,
		messageID,
		payload.OrderID,
		payload.CategoryID,
	); err != nil {
		return fmt.Errorf("view category products: %w", err)
	}

	return nil
}
