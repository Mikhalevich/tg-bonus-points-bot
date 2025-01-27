package messagesender

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (m *messageSender) SendOrderInvoice(
	ctx context.Context,
	chatID msginfo.ChatID,
	title string,
	description string,
	ord order.Order,
	rows ...button.InlineKeyboardButtonRow,
) error {
	if _, err := m.bot.SendInvoice(ctx, &bot.SendInvoiceParams{
		ChatID:        chatID.Int64(),
		Title:         title,
		Description:   description,
		Payload:       ord.ID.String(),
		ProviderToken: m.paymentToken,
		Currency:      "BYN",
		Prices:        convertProductsToLabeledPrices(ord.Products),
		ReplyMarkup:   makeButtonsMarkup(rows...),
	}); err != nil {
		return fmt.Errorf("send invoice: %w", err)
	}

	return nil
}

func convertProductsToLabeledPrices(products []order.OrderedProduct) []models.LabeledPrice {
	prices := make([]models.LabeledPrice, 0, len(products))

	for _, v := range products {
		prices = append(prices, models.LabeledPrice{
			Label:  fmt.Sprintf("%s x%d", v.Product.Title, v.Count),
			Amount: v.Count * v.Product.Price,
		})
	}

	return prices
}
