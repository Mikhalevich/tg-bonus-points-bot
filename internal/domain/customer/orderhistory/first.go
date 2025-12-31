package orderhistory

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/orderhistory/internal/page"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

func (o *OrderHistory) First(ctx context.Context, info msginfo.Info) error {
	twoPageOrders, err := o.repository.HistoryOrdersFirst(ctx, info.ChatID, o.pageSize+1)
	if err != nil {
		return fmt.Errorf("history orders: %w", err)
	}

	if len(twoPageOrders) == 0 {
		o.sendPlainText(ctx, info.ChatID, message.OrderNoOrdersFound())

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
		afterOrderIDBtn  = order.IDFromInt(0)
		beforeOrderIDBtn = page.CalculateOrderIDForNextPage(twoPageOrders, o.pageSize)
		onePageOrders    = page.TruncateOrdersToPageSize(twoPageOrders, o.pageSize)
		pageInfo         = page.Last(totalOrders, o.pageSize, page.DESC)
	)

	buttons, err := o.makeHistoryButtons(
		info.ChatID,
		afterOrderIDBtn,
		beforeOrderIDBtn,
	)
	if err != nil {
		return fmt.Errorf("make history buttons: %w", err)
	}

	o.editPlainText(
		ctx,
		info.ChatID,
		info.MessageID,
		formatHistoryOrders(onePageOrders, curr, pageInfo),
		buttons...,
	)

	return nil
}
