package customerorder

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

const (
	pageSize = 10
)

func (c *CustomerOrder) History(ctx context.Context, chatID msginfo.ChatID) error {
	orders, err := c.repository.HistoryOrders(ctx, chatID, pageSize)
	if err != nil {
		return fmt.Errorf("history orders: %w", err)
	}

	c.sender.SendTextMarkdown(ctx, chatID, formatShortOrders(orders, c.sender.EscapeMarkdown))

	return nil
}

func formatShortOrders(
	orders []order.ShortOrder,
	escaper func(string) string,
) string {
	if len(orders) == 0 {
		return message.OrderNoOrdersFound()
	}

	formattedOrders := make([]string, 0, len(orders))

	for _, v := range orders {
		formattedOrders = append(formattedOrders,
			fmt.Sprintf("created\\_at: *%s* status: *%s* price: *%d*",
				escaper(v.CreatedAt.Format(time.RFC3339)),
				v.Status.HumanReadable(),
				v.TotalPrice,
			),
		)
	}

	return strings.Join(formattedOrders, "\n")
}
