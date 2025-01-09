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

	info := msginfo.Info{
		ChatID:    msginfo.ChatIDFromInt(msg.ChatID),
		MessageID: msginfo.MessageIDFromInt(msg.MessageID),
	}

	switch btn.Operation {
	case button.OperationCancelOrderEditMessage:
		if err := t.cancelOrderEditMsg(ctx, info, btn); err != nil {
			return fmt.Errorf("cancel order edit msg: %w", err)
		}

	case button.OperationCancelOrderSendMessage:
		if err := t.cancelOrderSendMsg(ctx, info.ChatID, btn); err != nil {
			return fmt.Errorf("cancel order send msg: %w", err)
		}

	case button.OperationConfirmOrder:
		if err := t.confirmOrder(ctx, info, btn); err != nil {
			return fmt.Errorf("confirm order: %w", err)
		}

	case button.OperationViewCategory:
		if err := t.viewProducts(ctx, info, btn); err != nil {
			return fmt.Errorf("view products: %w", err)
		}

	case button.OperationProduct:
		return errors.New("not implemented")

	case button.OperationBackToOrder:
		if err := t.backToOrder(ctx, info, btn); err != nil {
			return fmt.Errorf("back to order: %w", err)
		}
	}

	return nil
}

func (t *TGHandler) cancelOrderEditMsg(ctx context.Context, info msginfo.Info, btn *button.Button) error {
	orderID, err := btn.OrderID()
	if err != nil {
		return fmt.Errorf("invalid order id: %w", err)
	}

	if err := t.orderProcessor.CancelOrderEditMessage(ctx, info, orderID); err != nil {
		return fmt.Errorf("cancel order: %w", err)
	}

	return nil
}

func (t *TGHandler) cancelOrderSendMsg(ctx context.Context, chatID msginfo.ChatID, btn *button.Button) error {
	orderID, err := btn.OrderID()
	if err != nil {
		return fmt.Errorf("invalid order id: %w", err)
	}

	if err := t.orderProcessor.CancelOrderSendMessage(ctx, chatID, orderID); err != nil {
		return fmt.Errorf("cancel order: %w", err)
	}

	return nil
}

func (t *TGHandler) confirmOrder(ctx context.Context, info msginfo.Info, btn *button.Button) error {
	orderID, err := btn.OrderID()
	if err != nil {
		return fmt.Errorf("invalid order id: %w", err)
	}

	if err := t.orderProcessor.ConfirmOrder(ctx, info, orderID); err != nil {
		return fmt.Errorf("confirm order: %w", err)
	}

	return nil
}

func (t *TGHandler) viewProducts(ctx context.Context, info msginfo.Info, btn *button.Button) error {
	payload, err := btn.ViewCategoryPayload()
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.orderProcessor.ViewCategoryProducts(
		ctx,
		info,
		payload.OrderID,
		payload.CategoryID,
	); err != nil {
		return fmt.Errorf("view category products: %w", err)
	}

	return nil
}

func (t *TGHandler) backToOrder(ctx context.Context, info msginfo.Info, btn *button.Button) error {
	orderID, err := btn.OrderID()
	if err != nil {
		return fmt.Errorf("invalid order id: %w", err)
	}

	if err := t.orderProcessor.RefreshOrder(ctx, info, orderID); err != nil {
		return fmt.Errorf("refresh order: %w", err)
	}

	return nil
}
