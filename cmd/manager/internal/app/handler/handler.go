package handler

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

type OrderProcessor interface {
	GetNextPendingOrderToProcess(ctx context.Context) (*order.Order, error)
	UpdateOrderStatus(ctx context.Context, id order.ID, status order.Status) error
}

type Handler struct {
	orderProcessor OrderProcessor
}

func New(orderProcessor OrderProcessor) *Handler {
	return &Handler{
		orderProcessor: orderProcessor,
	}
}
