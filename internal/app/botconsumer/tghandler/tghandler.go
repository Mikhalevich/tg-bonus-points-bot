package tghandler

import (
	"context"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

type OrderProcessor interface {
	MakeOrder(ctx context.Context, chatID msginfo.ChatID, messageID msginfo.MessageID) error
	GetActiveOrder(ctx context.Context, chatID msginfo.ChatID, messageID msginfo.MessageID) error
	CancelOrder(ctx context.Context, id order.ID) error
}

type TGHandler struct {
	orderProcessor OrderProcessor
}

func New(orderProcessor OrderProcessor) *TGHandler {
	return &TGHandler{
		orderProcessor: orderProcessor,
	}
}
