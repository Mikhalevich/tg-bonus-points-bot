package orderhistory

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (o *OrderHistory) History(ctx context.Context, chatID msginfo.ChatID) error {
	orders, err := o.repository.HistoryOrders(ctx, chatID, o.pageSize+1)
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

	buttons, err := o.makeHistoryButtons(ctx, chatID, orders)
	if err != nil {
		return fmt.Errorf("make history buttons: %w", err)
	}

	o.sender.SendTextMarkdown(ctx, chatID, formatHistoryOrders(orders, curr, o.sender.EscapeMarkdown), buttons...)

	return nil
}

func (o *OrderHistory) makeHistoryButtons(
	ctx context.Context,
	chatID msginfo.ChatID,
	orders []order.HistoryOrder,
) ([]button.InlineKeyboardButtonRow, error) {
	var buttons button.ButtonRow

	if len(orders) > o.pageSize {
		previousOrdersBtn, err := button.OrderHistoryPrevious(chatID, message.Previous(), orders[o.pageSize].ID)
		if err != nil {
			return nil, fmt.Errorf("previous history button: %w", err)
		}

		buttons = append(buttons, previousOrdersBtn)
	}

	inlineButtons, err := o.buttonRepository.SetButtonRows(ctx, buttons)
	if err != nil {
		return nil, fmt.Errorf("store buttons: %w", err)
	}

	return inlineButtons, nil
}

func formatHistoryOrders(
	orders []order.HistoryOrder,
	curr *currency.Currency,
	escaper func(string) string,
) string {
	formattedOrders := make([]string, 0, len(orders))

	for _, v := range orders {
		formattedOrders = append(formattedOrders,
			fmt.Sprintf("%d\\. *%s* *%s* *%s*",
				v.SerialNumber,
				escaper(v.CreatedAt.Format(time.RFC3339)),
				v.Status.HumanReadable(),
				curr.FormatPrice(v.TotalPrice),
			),
		)
	}

	return strings.Join(formattedOrders, "\n")
}
