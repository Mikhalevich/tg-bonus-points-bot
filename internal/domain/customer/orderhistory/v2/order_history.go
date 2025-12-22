package v2

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

type Repository interface {
	HistoryOrdersByOffset(ctx context.Context, chatID msginfo.ChatID, offset, limit int) ([]order.HistoryOrder, error)
	HistoryOrdersCount(ctx context.Context, chatID msginfo.ChatID) (int, error)
}

type CurrencyProvider interface {
	GetCurrencyByID(ctx context.Context, id currency.ID) (*currency.Currency, error)
}

type MessageSender interface {
	SendText(ctx context.Context, chatID msginfo.ChatID, text string,
		buttons ...button.InlineKeyboardButtonRow)
	SendTextMarkdown(ctx context.Context, chatID msginfo.ChatID, text string,
		buttons ...button.InlineKeyboardButtonRow)
	EditTextMessage(ctx context.Context, chatID msginfo.ChatID, messageID msginfo.MessageID,
		text string, rows ...button.InlineKeyboardButtonRow,
	)
}

type ButtonRowsSetter interface {
	SetButtonRows(ctx context.Context, rows ...button.ButtonRow) ([]button.InlineKeyboardButtonRow, error)
}

type OrderHistory struct {
	repository       Repository
	currencyProvider CurrencyProvider
	sender           MessageSender
	buttonRowsSetter ButtonRowsSetter
	pageSize         int
}

func New(
	repository Repository,
	currencyProvider CurrencyProvider,
	sender MessageSender,
	buttonRowsSetter ButtonRowsSetter,
	pageSize int,
) *OrderHistory {
	return &OrderHistory{
		repository:       repository,
		currencyProvider: currencyProvider,
		sender:           sender,
		buttonRowsSetter: buttonRowsSetter,
		pageSize:         pageSize,
	}
}
