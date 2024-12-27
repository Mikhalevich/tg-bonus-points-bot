package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/perror"
)

func (m *Manager) UpdateOrderStatus(ctx context.Context, id order.ID, status order.Status) error {
	if err := m.repository.UpdateOrderStatus(
		ctx,
		id,
		time.Now(),
		status,
		calculateLegalPreviousStatuses(status)...,
	); err != nil {
		if m.repository.IsNotUpdatedError(err) {
			return perror.NotFound("order with relevant status not found")
		}

		return fmt.Errorf("update order status: %w", err)
	}

	return nil
}

func calculateLegalPreviousStatuses(s order.Status) []order.Status {
	switch s {
	case order.StatusCreated:
		return nil
	case order.StatusInProgress:
		return []order.Status{order.StatusCreated}
	case order.StatusReady:
		return []order.Status{order.StatusInProgress}
	case order.StatusCompleted:
		return []order.Status{order.StatusCreated, order.StatusInProgress, order.StatusReady}
	case order.StatusCanceled:
		return []order.Status{order.StatusCreated, order.StatusInProgress, order.StatusReady}
	}

	return nil
}
