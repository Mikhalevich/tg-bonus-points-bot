package orderhistory

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/orderhistory/internal/page"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

func (o *OrderHistory) Show(ctx context.Context, chatID msginfo.ChatID) error {
	twoPageOrders, err := o.repository.HistoryOrdersFirst(ctx, chatID, o.pageSize+1)
	if err != nil {
		return fmt.Errorf("history orders: %w", err)
	}

	if len(twoPageOrders) == 0 {
		o.sendPlainText(ctx, chatID, message.OrderNoOrdersFound())

		return nil
	}

	totalOrders, err := o.repository.HistoryOrdersCount(ctx, chatID)
	if err != nil {
		return fmt.Errorf("history orders count: %w", err)
	}

	curr, err := o.currencyProvider.GetCurrencyByID(ctx, twoPageOrders[0].CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	var (
		afterOrderIDBtn  = order.IDFromInt(0)
		beforeOrderIDBtn = page.CalculateOrderIDForNextPage(twoPageOrders, o.pageSize)
		onePageOrders    = page.TruncateOrdersToPageSize(twoPageOrders, o.pageSize)
		pageInfo         = page.Last(totalOrders, o.pageSize, page.DESC)
	)

	buttons, err := o.makeHistoryButtons(
		chatID,
		afterOrderIDBtn,
		beforeOrderIDBtn,
	)
	if err != nil {
		return fmt.Errorf("make history buttons: %w", err)
	}

	o.sendPlainText(
		ctx,
		chatID,
		formatHistoryOrders(onePageOrders, curr, pageInfo),
		buttons...,
	)

	return nil
}

func (o *OrderHistory) makeHistoryButtons(
	chatID msginfo.ChatID,
	afterOrderID order.ID,
	beforeOrderID order.ID,
) ([]button.ButtonRow, error) {
	var buttons button.ButtonRow

	if afterOrderID.Int() > 0 {
		nextOrdersBtn, err := button.OrderHistoryByIDNext(chatID, message.HistoryNext(), afterOrderID)
		if err != nil {
			return nil, fmt.Errorf("previous history button: %w", err)
		}

		firstOrdersBtn := button.OrderHistoryByIDFirst(chatID, message.HistoryFirst())

		buttons = append(buttons, firstOrdersBtn, nextOrdersBtn)
	}

	if beforeOrderID.Int() > 0 {
		previousOrdersBtn, err := button.OrderHistoryByIDPrevious(chatID, message.HistoryPrevious(), beforeOrderID)
		if err != nil {
			return nil, fmt.Errorf("previous history button: %w", err)
		}

		lastOrdersBtn := button.OrderHistoryByIDLast(chatID, message.HistoryLast())

		buttons = append(buttons, previousOrdersBtn, lastOrdersBtn)
	}

	return []button.ButtonRow{buttons}, nil
}

func formatHistoryOrders(
	orders []order.HistoryOrder,
	curr *currency.Currency,
	pageInfo page.Page,
) string {
	formattedOrders := make([]string, 0, len(orders)+1)
	formattedOrders = append(formattedOrders, fmt.Sprintf("Page %d/%d", pageInfo.Number, pageInfo.Total))

	for _, ord := range orders {
		formattedOrders = append(formattedOrders,
			fmt.Sprintf("%d) %s, %s, %s",
				ord.SerialNumber,
				ord.CreatedAt.Format(time.RFC3339),
				ord.Status.HumanReadable(),
				curr.FormatPrice(ord.TotalPrice),
			),
		)
	}

	return strings.Join(formattedOrders, "\n")
}
