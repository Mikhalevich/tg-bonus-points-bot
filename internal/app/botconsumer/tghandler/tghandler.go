package tghandler

import (
	"context"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

type OrderProcessor interface {
	MakeOrder(ctx context.Context, chatID msginfo.ChatID, messageID msginfo.MessageID) error
	GetActiveOrder(ctx context.Context, chatID msginfo.ChatID, messageID msginfo.MessageID) error
	CancelOrder(ctx context.Context, chatID msginfo.ChatID, orderID order.ID) error
	ConfirmOrder(ctx context.Context, chatID msginfo.ChatID, orderID order.ID) error
	ViewCategoryProducts(ctx context.Context, chatID msginfo.ChatID, messageID msginfo.MessageID,
		categoryID product.ID) error
	GetButton(ctx context.Context, id button.ID) (*button.Button, error)
}

type TGHandler struct {
	orderProcessor OrderProcessor
}

func New(orderProcessor OrderProcessor) *TGHandler {
	return &TGHandler{
		orderProcessor: orderProcessor,
	}
}
