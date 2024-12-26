package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/perror"
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

func convertErrorToHTTPCode(err error) int {
	var perr *perror.Error
	if !errors.As(err, &perr) {
		return http.StatusInternalServerError
	}

	if perr.Type == perror.TypeNotFound {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}
