package orderhistory

import (
	"context"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

type CurrencyProvider interface {
	GetCurrencyByID(ctx context.Context, id currency.ID) (*currency.Currency, error)
}

type Repository interface {
	HistoryOrdersCount(ctx context.Context, chatID msginfo.ChatID) (int, error)
	HistoryOrdersFirst(ctx context.Context, chatID msginfo.ChatID, size int) ([]order.HistoryOrder, error)
	HistoryOrdersLast(ctx context.Context, chatID msginfo.ChatID, size int) ([]order.HistoryOrder, error)
	HistoryOrdersBeforeID(ctx context.Context, chatID msginfo.ChatID, id order.ID, size int) ([]order.HistoryOrder, error)
	HistoryOrdersAfterID(ctx context.Context, chatID msginfo.ChatID, id order.ID, size int) ([]order.HistoryOrder, error)
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
	currencyProvider CurrencyProvider
	repository       Repository
	sender           MessageSender
	buttonSetter     ButtonRowsSetter
	pageSize         int
}

func New(
	currencyProvider CurrencyProvider,
	repository Repository,
	sender port.MessageSender,
	buttonSetter ButtonRowsSetter,
	pageSize int,
) *OrderHistory {
	return &OrderHistory{
		currencyProvider: currencyProvider,
		repository:       repository,
		sender:           sender,
		buttonSetter:     buttonSetter,
		pageSize:         pageSize,
	}
}
