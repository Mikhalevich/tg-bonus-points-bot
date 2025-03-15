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
	twoPageOrders, err := o.repository.HistoryOrders(ctx, chatID, o.pageSize+1)
	if err != nil {
		return fmt.Errorf("history orders: %w", err)
	}

	if len(twoPageOrders) == 0 {
		o.sender.SendTextMarkdown(ctx, chatID, message.OrderNoOrdersFound())
		return nil
	}

	curr, err := o.repository.GetCurrencyByID(ctx, twoPageOrders[0].CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	buttons, err := o.makeHistoryButtons(
		ctx,
		chatID,
		0,
		calculateOrderIDForPreviousPage(twoPageOrders, o.pageSize),
	)
	if err != nil {
		return fmt.Errorf("make history buttons: %w", err)
	}

	o.sender.SendText(
		ctx,
		chatID,
		formatHistoryOrders(truncateOrdersToPageSize(twoPageOrders, o.pageSize), curr),
		buttons...,
	)

	return nil
}

func calculateOrderIDForPreviousPage(twoPageOrders []order.HistoryOrder, pageSize int) order.ID {
	if len(twoPageOrders) > pageSize {
		return twoPageOrders[pageSize-1].ID
	}

	return 0
}

func calculateOrderIDForNextPage(twoPageOrders []order.HistoryOrder, pageSize int) order.ID {
	if len(twoPageOrders) > pageSize {
		return twoPageOrders[0].ID
	}

	return 0
}

func truncateOrdersToPageSize(orders []order.HistoryOrder, pageSize int) []order.HistoryOrder {
	if len(orders) > pageSize {
		return orders[:pageSize]
	}

	return orders
}

func (o *OrderHistory) makeHistoryButtons(
	ctx context.Context,
	chatID msginfo.ChatID,
	afterOrderID order.ID,
	beforeOrderID order.ID,
) ([]button.InlineKeyboardButtonRow, error) {
	var buttons button.ButtonRow

	if afterOrderID.Int() > 0 {
		nextOrdersBtn, err := button.OrderHistoryNext(chatID, message.Next(), afterOrderID)
		if err != nil {
			return nil, fmt.Errorf("previous history button: %w", err)
		}

		buttons = append(buttons, nextOrdersBtn)
	}

	if beforeOrderID.Int() > 0 {
		previousOrdersBtn, err := button.OrderHistoryPrevious(chatID, message.Previous(), beforeOrderID)
		if err != nil {
			return nil, fmt.Errorf("previous history button: %w", err)
		}

		buttons = append(buttons, previousOrdersBtn)
	}

	if len(buttons) == 0 {
		return nil, nil
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
) string {
	formattedOrders := make([]string, 0, len(orders))

	for _, v := range orders {
		formattedOrders = append(formattedOrders,
			fmt.Sprintf("%d) %s, %s, %s",
				v.SerialNumber,
				v.CreatedAt.Format(time.RFC3339),
				v.Status.HumanReadable(),
				curr.FormatPrice(v.TotalPrice),
			),
		)
	}

	return strings.Join(formattedOrders, "\n")
}
