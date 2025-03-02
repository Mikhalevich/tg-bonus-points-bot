package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/internal/tgbot"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (t *TGHandler) GetOrderbyID(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	orderID, err := order.IDFromString(msg.Text)
	if err != nil {
		return fmt.Errorf("invalid order id: %w", err)
	}

	if err := t.orderProcessor.GetOrderByID(
		ctx,
		msginfo.ChatIDFromInt(msg.ChatID),
		orderID,
	); err != nil {
		return fmt.Errorf("get active order: %w", err)
	}

	return nil
}
