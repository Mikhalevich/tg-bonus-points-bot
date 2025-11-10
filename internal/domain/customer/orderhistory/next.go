package orderhistory

import (
	"context"
	"fmt"
	"slices"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/customer/orderhistory/internal/page"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (o *OrderHistory) Next(
	ctx context.Context,
	info msginfo.Info,
	afterOrderID order.ID,
) error {
	twoPageOrders, err := o.repository.HistoryOrdersAfterID(ctx, info.ChatID, afterOrderID, o.pageSize+1)
	if err != nil {
		return fmt.Errorf("history orders: %w", err)
	}

	if len(twoPageOrders) == 0 {
		o.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, message.OrderNoOrdersFound())

		return nil
	}

	totalOrders, err := o.repository.HistoryOrdersCount(ctx, info.ChatID)
	if err != nil {
		return fmt.Errorf("history orders count: %w", err)
	}

	curr, err := o.currencyProvider.GetCurrencyByID(ctx, twoPageOrders[0].CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	var (
		afterOrderIDBtn  = page.CalculateOrderIDForNextPage(twoPageOrders, o.pageSize)
		onePageOrders    = page.TruncateOrdersToPageSize(twoPageOrders, o.pageSize)
		beforeOrderIDBtn = onePageOrders[0].ID
		pageInfo         = page.Current(onePageOrders, totalOrders, page.BottomToTop, o.pageSize, page.DESC)
	)

	buttons, err := o.makeHistoryButtons(
		ctx,
		info.ChatID,
		afterOrderIDBtn,
		beforeOrderIDBtn,
	)
	if err != nil {
		return fmt.Errorf("make history buttons: %w", err)
	}

	slices.Reverse(onePageOrders)

	o.sender.EditTextMessage(
		ctx,
		info.ChatID,
		info.MessageID,
		formatHistoryOrders(onePageOrders, curr, pageInfo),
		buttons...,
	)

	return nil
}
