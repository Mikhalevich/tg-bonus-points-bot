package handler

import (
	"context"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

type Manager interface {
	GetNextPendingOrderToProcess(ctx context.Context) (*order.Order, error)
}

type Handler struct {
	manager Manager
}

func New(manager Manager) *Handler {
	return &Handler{
		manager: manager,
	}
}
