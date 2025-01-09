package tghandler

import (
	"context"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

type OrderProcessor interface {
	MakeOrder(ctx context.Context, info msginfo.Info) error
	GetActiveOrder(ctx context.Context, info msginfo.Info) error
	CancelOrderSendMessage(ctx context.Context, chatID msginfo.ChatID, orderID order.ID) error
	CancelOrderEditMessage(ctx context.Context, info msginfo.Info, orderID order.ID) error
	ConfirmOrder(ctx context.Context, info msginfo.Info, orderID order.ID) error
	ViewCategoryProducts(ctx context.Context, info msginfo.Info, orderID order.ID, categoryID product.ID) error
	RefreshOrder(ctx context.Context, info msginfo.Info, orderID order.ID) error
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
