package orderhistory

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/customer/orderhistory/internal/page"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (o *OrderHistory) Show(ctx context.Context, chatID msginfo.ChatID) error {
	twoPageOrders, err := o.repository.HistoryOrdersFirst(ctx, chatID, o.pageSize+1)
	if err != nil {
		return fmt.Errorf("history orders: %w", err)
	}

	if len(twoPageOrders) == 0 {
		o.sender.SendTextMarkdown(ctx, chatID, message.OrderNoOrdersFound())
		return nil
	}

	totalOrders, err := o.repository.HistoryOrdersCount(ctx, chatID)
	if err != nil {
		return fmt.Errorf("history orders count: %w", err)
	}

	curr, err := o.repository.GetCurrencyByID(ctx, twoPageOrders[0].CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	var (
		afterOrderIDBtn  = order.IDFromInt(0)
		beforeOrderIDBtn = page.CalculateOrderIDForNextPage(twoPageOrders, o.pageSize)
		onePageOrders    = page.TruncateOrdersToPageSize(twoPageOrders, o.pageSize)
		pageInfo         = page.Last(onePageOrders, totalOrders, o.pageSize)
	)

	buttons, err := o.makeHistoryButtons(
		ctx,
		chatID,
		afterOrderIDBtn,
		beforeOrderIDBtn,
	)
	if err != nil {
		return fmt.Errorf("make history buttons: %w", err)
	}

	o.sender.SendText(
		ctx,
		chatID,
		formatHistoryOrders(onePageOrders, curr, pageInfo),
		buttons...,
	)

	return nil
}

func (o *OrderHistory) makeHistoryButtons(
	ctx context.Context,
	chatID msginfo.ChatID,
	afterOrderID order.ID,
	beforeOrderID order.ID,
) ([]button.InlineKeyboardButtonRow, error) {
	var buttons button.ButtonRow

	if afterOrderID.Int() > 0 {
		nextOrdersBtn, err := button.OrderHistoryNext(chatID, message.HistoryNext(), afterOrderID)
		if err != nil {
			return nil, fmt.Errorf("previous history button: %w", err)
		}

		firstOrdersBtn := button.OrderHistoryFirst(chatID, message.HistoryFirst())

		buttons = append(buttons, firstOrdersBtn, nextOrdersBtn)
	}

	if beforeOrderID.Int() > 0 {
		previousOrdersBtn, err := button.OrderHistoryPrevious(chatID, message.HistoryPrevious(), beforeOrderID)
		if err != nil {
			return nil, fmt.Errorf("previous history button: %w", err)
		}

		lastOrdersBtn := button.OrderHistoryLast(chatID, message.HistoryLast())

		buttons = append(buttons, previousOrdersBtn, lastOrdersBtn)
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
	pageInfo page.Page,
) string {
	formattedOrders := make([]string, 0, len(orders)+1)
	formattedOrders = append(formattedOrders, fmt.Sprintf("Page %d/%d", pageInfo.Number, pageInfo.Total))

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
