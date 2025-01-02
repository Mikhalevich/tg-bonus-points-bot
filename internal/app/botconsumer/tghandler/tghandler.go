package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/internal/tgbot"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

type OrderProcessor interface {
	MakeOrder(ctx context.Context, chatID msginfo.ChatID, messageID msginfo.MessageID) error
	GetActiveOrder(ctx context.Context, chatID msginfo.ChatID, messageID msginfo.MessageID) error
	CancelOrder(ctx context.Context, id order.ID) error
	GetButton(ctx context.Context, id button.ID) (*button.Button, error)
	TestProducts(ctx context.Context, chatID msginfo.ChatID) error
}

type TGHandler struct {
	orderProcessor OrderProcessor
}

func New(orderProcessor OrderProcessor) *TGHandler {
	return &TGHandler{
		orderProcessor: orderProcessor,
	}
}

func (t *TGHandler) TestProducts(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if err := t.orderProcessor.TestProducts(ctx, msginfo.ChatID(msg.ChatID)); err != nil {
		return fmt.Errorf("test products: %w", err)
	}
	return nil
}
