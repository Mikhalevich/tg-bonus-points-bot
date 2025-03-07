package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/internal/tgbot"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/perror"
)

func (t *TGHandler) DefaultCallbackQuery(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	btn, err := t.buttonProvider.GetButton(ctx, button.IDFromString(msg.Data))
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

	handler, ok := t.cbHandlers[btn.Operation]
	if !ok {
		return fmt.Errorf("operation %s is not implented", btn.Operation)
	}

	if err := handler(ctx, info, *btn); err != nil {
		return fmt.Errorf("cb handler operation %s failure: %w", btn.Operation, err)
	}

	return nil
}

func (t *TGHandler) cancelOrder(ctx context.Context, info msginfo.Info, btn button.Button) error {
	payload, err := button.GetPayload[button.CancelOrderPayload](btn)
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.actionProcessor.Cancel(
		ctx,
		info.ChatID,
		info.MessageID,
		payload.OrderID,
		payload.IsTextMessage,
	); err != nil {
		return fmt.Errorf("cancel order: %w", err)
	}

	return nil
}

func (t *TGHandler) confirmCart(ctx context.Context, info msginfo.Info, btn button.Button) error {
	payload, err := button.GetPayload[button.CartConfirmPayload](btn)
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.cartProcessor.Confirm(ctx, info, payload.CartID, payload.CurrencyID); err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	return nil
}

func (t *TGHandler) cancelCart(ctx context.Context, info msginfo.Info, btn button.Button) error {
	payload, err := button.GetPayload[button.CartCancelPayload](btn)
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.cartProcessor.Cancel(ctx, info, payload.CartID); err != nil {
		return fmt.Errorf("cart cancel: %w", err)
	}

	return nil
}

func (t *TGHandler) viewCategoryProducts(ctx context.Context, info msginfo.Info, btn button.Button) error {
	payload, err := button.GetPayload[button.CartViewCategoryProductsPayload](btn)
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.cartProcessor.ViewCategoryProducts(
		ctx,
		info,
		payload.CartID,
		payload.CategoryID,
		payload.CurrencyID,
	); err != nil {
		return fmt.Errorf("view category products: %w", err)
	}

	return nil
}

func (t *TGHandler) viewCategories(ctx context.Context, info msginfo.Info, btn button.Button) error {
	payload, err := button.GetPayload[button.CartViewCategoriesPayload](btn)
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.cartProcessor.ViewCategories(ctx, info, payload.CartID, payload.CurrencyID); err != nil {
		return fmt.Errorf("cart view categories: %w", err)
	}

	return nil
}

func (t *TGHandler) addProduct(ctx context.Context, info msginfo.Info, btn button.Button) error {
	payload, err := button.GetPayload[button.CartAddProductPayload](btn)
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.cartProcessor.AddProduct(
		ctx,
		info,
		payload.CartID,
		payload.CategoryID,
		payload.ProductID,
		payload.CurrencyID,
	); err != nil {
		return fmt.Errorf("cart add product: %w", err)
	}

	return nil
}
