package orderhistory

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (o *OrderHistory) History(ctx context.Context, chatID msginfo.ChatID) error {
	orders, err := o.repository.HistoryOrders(ctx, chatID, o.pageSize)
	if err != nil {
		return fmt.Errorf("history orders: %w", err)
	}

	if len(orders) == 0 {
		o.sender.SendTextMarkdown(ctx, chatID, message.OrderNoOrdersFound())
		return nil
	}

	curr, err := o.repository.GetCurrencyByID(ctx, orders[0].CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	o.sender.SendTextMarkdown(ctx, chatID, formatHistoryOrders(orders, curr, o.sender.EscapeMarkdown))

	return nil
}

func formatHistoryOrders(
	orders []order.HistoryOrder,
	curr *currency.Currency,
	escaper func(string) string,
) string {
	if len(orders) == 0 {
		return message.OrderNoOrdersFound()
	}

	formattedOrders := make([]string, 0, len(orders))

	for _, v := range orders {
		formattedOrders = append(formattedOrders,
			fmt.Sprintf("%d\\. created\\_at: *%s* status: *%s* price: *%s*",
				v.SerialNumber,
				escaper(v.CreatedAt.Format(time.RFC3339)),
				v.Status.HumanReadable(),
				curr.FormatPrice(v.TotalPrice),
			),
		)
	}

	return strings.Join(formattedOrders, "\n")
}
