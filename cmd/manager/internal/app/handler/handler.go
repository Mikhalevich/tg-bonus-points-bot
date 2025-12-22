package handler

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

type Manager interface {
	GetNextPendingOrderToProcess(ctx context.Context) (*order.Order, error)
	UpdateOrderStatus(ctx context.Context, id order.ID, status order.Status) error
	GetProducts(ctx context.Context, filter product.Filter) ([]product.CategoryProducts, error)
}

type Handler struct {
	manager Manager
}

func New(manager Manager) *Handler {
	return &Handler{
		manager: manager,
	}
}
