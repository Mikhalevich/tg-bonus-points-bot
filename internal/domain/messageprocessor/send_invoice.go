package messageprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func (m *MessageProcessor) SendInvoice(
	ctx context.Context,
	chatID msginfo.ChatID,
	title string,
	description string,
	ord *order.Order,
	productsInfo map[product.ProductID]product.Product,
	currency string,
	rows ...button.ButtonRow,
) error {
	inlineButtons, err := m.SetButtonRows(ctx, rows...)
	if err != nil {
		return fmt.Errorf("set button rows: %w", err)
	}

	if err := m.sender.SendOrderInvoice(
		ctx,
		chatID,
		title,
		description,
		ord,
		productsInfo,
		currency,
		inlineButtons...,
	); err != nil {
		return fmt.Errorf("send order invoice: %w", err)
	}

	return nil
}
