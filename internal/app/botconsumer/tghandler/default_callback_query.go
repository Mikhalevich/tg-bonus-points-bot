package tghandler

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/internal/tgbot"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/perror"
)

func (t *TGHandler) DefaultCallbackQuery(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	btn, err := t.orderProcessor.GetButton(ctx, button.IDFromString(msg.Data))
	if err != nil {
		if perror.IsType(err, perror.TypeNotFound) {
			sender.SendMessage(ctx, msg.ChatID, "Action already executed or expired")
			return nil
		}

		return fmt.Errorf("get button: %w", err)
	}

	if btn.Operation == button.OperationCancelOrder {
		orderID, err := btn.OrderID()
		if err != nil {
			return errors.New("invalid order id")
		}

		if err := t.orderProcessor.CancelOrder(ctx, orderID); err != nil {
			return fmt.Errorf("cancel order: %w", err)
		}
	}

	return nil
}
