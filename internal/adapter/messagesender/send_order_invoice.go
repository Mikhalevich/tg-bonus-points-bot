package messagesender

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (m *messageSender) SendOrderInvoice(
	ctx context.Context,
	chatID msginfo.ChatID,
	title string,
	description string,
	ord *order.Order,
	productsInfo map[product.ProductID]product.Product,
	currency string,
	rows ...button.InlineKeyboardButtonRow,
) error {
	if _, err := m.bot.SendInvoice(ctx, &bot.SendInvoiceParams{
		ChatID:        chatID.Int64(),
		Title:         title,
		Description:   description,
		Payload:       ord.ID.String(),
		ProviderToken: m.paymentToken,
		Currency:      currency,
		Prices:        makeLabeledPrices(ord.Products, productsInfo),
		ReplyMarkup:   makeButtonsMarkup(rows...),
	}); err != nil {
		return fmt.Errorf("send invoice: %w", err)
	}

	return nil
}

func makeLabeledPrices(
	orderedProducts []order.OrderedProduct,
	productsInfo map[product.ProductID]product.Product,
) []models.LabeledPrice {
	prices := make([]models.LabeledPrice, 0, len(orderedProducts))

	for _, v := range orderedProducts {
		prices = append(prices, models.LabeledPrice{
			Label:  fmt.Sprintf("%s x%d", productsInfo[v.ProductID].Title, v.Count),
			Amount: v.Count * v.Price,
		})
	}

	return prices
}
