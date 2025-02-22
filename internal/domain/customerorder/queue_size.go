package customerorder

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (c *CustomerOrder) QueueSize(ctx context.Context, info msginfo.Info) error {
	count, err := c.repository.GetOrdersCountByStatus(ctx, order.StatusConfirmed, order.StatusInProgress)
	if err != nil {
		return fmt.Errorf("get orders count by status: %w", err)
	}

	c.sender.ReplyTextMarkdown(ctx, info.ChatID, info.MessageID, fmt.Sprintf("*%d*", count))

	return nil
}
