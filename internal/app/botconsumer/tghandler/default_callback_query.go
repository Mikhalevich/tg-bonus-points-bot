package tghandler

import (
	"context"
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
	case button.OperationOrderCancel:
		if err := t.cancelOrder(ctx, info.ChatID, btn); err != nil {
			return fmt.Errorf("cancel order edit msg: %w", err)
		}

	case button.OperationCartCancel:
		return fmt.Errorf("not implemented")

	case button.OperationCartConfirm:
		if err := t.confirmCart(ctx, info, btn); err != nil {
			return fmt.Errorf("confirm order: %w", err)
		}

	case button.OperationCartViewCategories:
		if err := t.viewCategories(ctx, info, btn); err != nil {
			return fmt.Errorf("view categories: %w", err)
		}

	case button.OperationCartViewCategoryProducts:
		if err := t.viewCategoryProducts(ctx, info, btn); err != nil {
			return fmt.Errorf("view products: %w", err)
		}

	case button.OperationCartAddProduct:
		if err := t.addProduct(ctx, info, btn); err != nil {
			return fmt.Errorf("add product: %w", err)
		}
	}

	return nil
}

func (t *TGHandler) cancelOrder(ctx context.Context, chatID msginfo.ChatID, btn *button.Button) error {
	orderID, err := btn.OrderID()
	if err != nil {
		return fmt.Errorf("invalid order id: %w", err)
	}

	if err := t.orderProcessor.CancelOrder(ctx, chatID, orderID); err != nil {
		return fmt.Errorf("cancel order: %w", err)
	}

	return nil
}

func (t *TGHandler) confirmCart(ctx context.Context, info msginfo.Info, btn *button.Button) error {
	payload, err := btn.CartConfirmPayload()
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.orderProcessor.CreateOrder(ctx, info, payload.CartID); err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	return nil
}

func (t *TGHandler) viewCategoryProducts(ctx context.Context, info msginfo.Info, btn *button.Button) error {
	payload, err := btn.CartCategoryProductsPayload()
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.orderProcessor.CartViewCategoryProducts(ctx, info, payload.CartID, payload.CategoryID); err != nil {
		return fmt.Errorf("view category products: %w", err)
	}

	return nil
}

func (t *TGHandler) viewCategories(ctx context.Context, info msginfo.Info, btn *button.Button) error {
	payload, err := btn.CartViewCategoriesPayload()
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.orderProcessor.CartViewCategories(ctx, info, payload.CartID); err != nil {
		return fmt.Errorf("cart view categories: %w", err)
	}

	return nil
}

func (t *TGHandler) addProduct(ctx context.Context, info msginfo.Info, btn *button.Button) error {
	payload, err := btn.AddProductPayload()
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.orderProcessor.CartAddProduct(
		ctx,
		info,
		payload.CartID,
		payload.CategoryID,
		payload.ProductID,
	); err != nil {
		return fmt.Errorf("cart add product: %w", err)
	}

	return nil
}
