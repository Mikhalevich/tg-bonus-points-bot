package v2

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

type SendMessageStrategy int

const (
	SendMessage SendMessageStrategy = iota + 1
	EditMessage
)

func (oh *OrderHistory) Page(
	ctx context.Context,
	info msginfo.Info,
	pageNumber int,
) error {
	if err := oh.loadPageByNumber(ctx, info, pageNumber, EditMessage); err != nil {
		return fmt.Errorf("load page by number: %w", err)
	}

	return nil
}

func (oh *OrderHistory) loadPageByNumber(
	ctx context.Context,
	info msginfo.Info,
	pageNumber int,
	sendStrategy SendMessageStrategy,
) error {
	ordersCount, err := oh.repository.HistoryOrdersCount(ctx, info.ChatID)
	if err != nil {
		return fmt.Errorf("history orders count: %w", err)
	}

	if ordersCount == 0 {
		oh.sendPlainText(ctx, info.ChatID, message.OrderNoOrdersFound())

		return nil
	}

	pagesCount := calculatePageCount(ordersCount, oh.pageSize)

	if pageNumber > pagesCount {
		return fmt.Errorf("invalid page number: %d, pages total: %d", pageNumber, pagesCount)
	}

	if err := oh.loadPageByPageInfo(
		ctx,
		info,
		page.Page{
			Number: pageNumber,
			Total:  pagesCount,
		},
		sendStrategy,
	); err != nil {
		return fmt.Errorf("load page by page info: %w", err)
	}

	return nil
}

func (oh *OrderHistory) loadPageByPageInfo(
	ctx context.Context,
	info msginfo.Info,
	pageInfo page.Page,
	sendStrategy SendMessageStrategy,
) error {
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
		oh.sendPlainText(ctx, info.ChatID, message.OrderNoOrdersFound())

		return nil
	}

	curr, err := oh.currencyProvider.GetCurrencyByID(ctx, historyOrders[0].CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	buttons, err := oh.makeHistoryButtons(
		info.ChatID,
		pageInfo,
	)
	if err != nil {
		return fmt.Errorf("make history buttons: %w", err)
	}

	switch sendStrategy {
	case SendMessage:
		oh.sendPlainText(
			ctx,
			info.ChatID,
			formatHistoryOrders(historyOrders, curr, pageInfo),
			buttons...,
		)

	case EditMessage:
		oh.editPlainText(
			ctx,
			info.ChatID,
			info.MessageID,
			formatHistoryOrders(historyOrders, curr, pageInfo),
			buttons...,
		)

	default:
		return fmt.Errorf("unknown send strategy: %v", sendStrategy)
	}

	return nil
}

func calculatePageOffset(pageNumber, pageSize int) int {
	if pageNumber <= 1 {
		return 0
	}

	return (pageNumber - 1) * pageSize
}

func calculatePageCount(count, pageSize int) int {
	var (
		fullPages    = count / pageSize
		lastPageSize = count % pageSize
	)

	if lastPageSize > 0 {
		return fullPages + 1
	}

	return fullPages
}

func formatHistoryOrders(
	orders []order.HistoryOrder,
	curr *currency.Currency,
	currentPage page.Page,
) string {
	formattedOrders := make([]string, 0, len(orders)+1)
	formattedOrders = append(formattedOrders, fmt.Sprintf("Page %d/%d", currentPage.Number, currentPage.Total))

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

func (oh *OrderHistory) makeHistoryButtons(
	chatID msginfo.ChatID,
	currentPage page.Page,
) ([]button.ButtonRow, error) {
	var buttons button.ButtonRow

	if currentPage.HasPrevious() {
		nextBtn, err := button.OrderHistoryByPage(chatID, message.HistoryNext(), currentPage.Previous())
		if err != nil {
			return nil, fmt.Errorf("next history button: %w", err)
		}

		firstBtn := button.OrderHistoryByPageFirst(chatID, message.HistoryFirst())

		buttons = append(buttons, firstBtn, nextBtn)
	}

	if currentPage.HasNext() {
		previousBtn, err := button.OrderHistoryByPage(chatID, message.HistoryPrevious(), currentPage.Next())
		if err != nil {
			return nil, fmt.Errorf("previous history button: %w", err)
		}

		lastBtn := button.OrderHistoryByPageLast(chatID, message.HistoryLast())

		buttons = append(buttons, previousBtn, lastBtn)
	}

	return []button.ButtonRow{buttons}, nil
}
