package orderhistory

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (o *OrderHistory) First(ctx context.Context, info msginfo.Info) error {
	twoPageOrders, err := o.repository.HistoryOrdersFirst(ctx, info.ChatID, o.pageSize+1)
	if err != nil {
		return fmt.Errorf("history orders: %w", err)
	}

	if len(twoPageOrders) == 0 {
		o.sender.SendTextMarkdown(ctx, info.ChatID, message.OrderNoOrdersFound())
		return nil
	}

	curr, err := o.repository.GetCurrencyByID(ctx, twoPageOrders[0].CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	var (
		afterOrderIDBtn  = order.IDFromInt(0)
		beforeOrderIDBtn = calculateOrderIDForNextPage(twoPageOrders, o.pageSize)
		onePageOrders    = truncateOrdersToPageSize(twoPageOrders, o.pageSize)
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

	o.sender.EditTextMessage(
		ctx,
		info.ChatID,
		info.MessageID,
		formatHistoryOrders(onePageOrders, curr),
		buttons...,
	)

	return nil
}
