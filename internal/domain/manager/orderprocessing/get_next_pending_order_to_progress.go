package orderprocessing

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (o *OrderProcessing) GetNextPendingOrderToProcess(ctx context.Context) (*order.Order, error) {
	order, err := o.repository.UpdateOrderStatusForMinID(
		ctx,
		o.timeProvider.Now(),
		order.StatusInProgress,
		order.StatusConfirmed,
	)
	if err != nil {
		if o.repository.IsNotFoundError(err) {
			return nil, perror.NotFound("no pending orders")
		}

		return nil, fmt.Errorf("update next order status: %w", err)
	}

	o.sendMarkdown(ctx, order.ChatID,
		o.makeChangedOrderStatusMarkdownMsg(order.Status))

	return order, nil
}

func (o *OrderProcessing) makeChangedOrderStatusMarkdownMsg(s order.Status) string {
	return fmt.Sprintf("your order status changed to *%s*", o.customerSender.EscapeMarkdown(s.HumanReadable()))
}
