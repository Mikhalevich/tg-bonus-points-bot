package v2

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/customer/orderhistory/internal/page"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

func (oh *OrderHistory) Last(ctx context.Context, info msginfo.Info) error {
	ordersCount, err := oh.repository.HistoryOrdersCount(ctx, info.ChatID)
	if err != nil {
		return fmt.Errorf("history orders count: %w", err)
	}

	if ordersCount == 0 {
		oh.sender.SendText(ctx, info.ChatID, message.OrderNoOrdersFound())

		return nil
	}

	pagesCount := calculatePageCount(ordersCount, oh.pageSize)

	if err := oh.loadPageByPageInfo(
		ctx,
		info,
		page.Page{
			Number: pagesCount,
			Total:  pagesCount,
		},
		EditMessage,
	); err != nil {
		return fmt.Errorf("load page by page info: %w", err)
	}

	return nil
}
