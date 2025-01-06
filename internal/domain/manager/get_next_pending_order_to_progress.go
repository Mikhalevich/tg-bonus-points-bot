package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/perror"
)

func (m *Manager) GetNextPendingOrderToProcess(ctx context.Context) (*order.Order, error) {
	order, err := m.repository.UpdateOrderStatusForMinID(ctx, time.Now(), order.StatusInProgress, order.StatusConfirmed)
	if err != nil {
		if m.repository.IsNotFoundError(err) {
			return nil, perror.NotFound("no pending orders")
		}

		return nil, fmt.Errorf("update next order status: %w", err)
	}

	m.customerSender.SendTextMarkdown(ctx, order.ChatID,
		m.makeChangedOrderStatusMarkdownMsg(order.Status))

	return order, nil
}

func (m *Manager) makeChangedOrderStatusMarkdownMsg(s order.Status) string {
	return fmt.Sprintf("your order status changed to *%s*", m.customerSender.EscapeMarkdown(s.HumanReadable()))
}
