package v2

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/customer/orderhistory/internal/page"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

func (oh *OrderHistory) Show(ctx context.Context, info msginfo.Info) error {
	ordersCount, err := oh.repository.HistoryOrdersCount(ctx, info.ChatID)
	if err != nil {
		return fmt.Errorf("history orders count: %w", err)
	}

	if ordersCount == 0 {
		oh.sender.SendText(ctx, info.ChatID, message.OrderNoOrdersFound())

		return nil
	}

	pagesCount := calculatePageCount(ordersCount, oh.pageSize)

	pageInfo := page.Page{
		Number: 1,
		Total:  pagesCount,
	}

	historyOrders, err := oh.repository.HistoryOrdersByOffset(
		ctx,
		info.ChatID,
		calculatePageOffset(pageInfo.Number, oh.pageSize),
		oh.pageSize,
	)
	if err != nil {
		return fmt.Errorf("history orders by offset: %w", err)
	}

	if len(historyOrders) == 0 {
		oh.sender.SendText(ctx, info.ChatID, message.OrderNoOrdersFound())

		return nil
	}

	curr, err := oh.currencyProvider.GetCurrencyByID(ctx, historyOrders[0].CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	buttons, err := oh.makeHistoryButtons(
		ctx,
		info.ChatID,
		pageInfo,
	)
	if err != nil {
		return fmt.Errorf("make history buttons: %w", err)
	}

	oh.sender.SendText(
		ctx,
		info.ChatID,
		formatHistoryOrders(historyOrders, curr, pageInfo),
		buttons...,
	)

	return nil
}
