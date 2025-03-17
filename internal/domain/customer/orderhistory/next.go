package orderhistory

import (
	"context"
	"fmt"

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

	curr, err := o.repository.GetCurrencyByID(ctx, twoPageOrders[0].CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	var (
		afterOrderIDBtn  = calculateOrderIDForNextPage(twoPageOrders, o.pageSize)
		onePageOrders    = truncateOrdersToPageSizeLeft(twoPageOrders, o.pageSize)
		beforeOrderIDBtn = onePageOrders[len(onePageOrders)-1].ID
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
