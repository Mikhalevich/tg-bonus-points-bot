package orderaction

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (o *OrderAction) QueueSize(ctx context.Context, info msginfo.Info) error {
	count, err := o.repository.GetOrdersCountByStatus(ctx, order.StatusConfirmed, order.StatusInProgress)
	if err != nil {
		return fmt.Errorf("get orders count by status: %w", err)
	}

	o.sender.ReplyTextMarkdown(ctx, info.ChatID, info.MessageID, fmt.Sprintf("*%d*", count))

	return nil
}
