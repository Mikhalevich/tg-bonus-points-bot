package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/perror"
)

func (m *Manager) GetNextPendingOrderToProcess(ctx context.Context) (*order.Order, error) {
	order, err := m.repository.UpdateOrderStatusForMinID(ctx, time.Now(), order.StatusCreated, order.StatusInProgress)
	if err != nil {
		if m.repository.IsNotFoundError(err) {
			return nil, perror.NotFound("no pending orders")
		}

		return nil, fmt.Errorf("update next order status: %w", err)
	}

	return order, nil
}
