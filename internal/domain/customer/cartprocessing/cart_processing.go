package cartprocessing

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
)

type ButtonRepositoryWriter interface {
	SetButton(ctx context.Context, b button.Button) (button.InlineKeyboardButton, error)
	SetButtonRows(ctx context.Context, rows ...button.ButtonRow) ([]button.InlineKeyboardButtonRow, error)
}

type TimeProvider interface {
	Now() time.Time
}

type CartProcessing struct {
	storeID          store.ID
	repository       port.CustomerCartRepository
	storeInfo        port.StoreInfo
	cart             port.Cart
	sender           port.MessageSender
	timeProvider     TimeProvider
	buttonRepository ButtonRepositoryWriter
}

func New(
	storeID int,
	repository port.CustomerCartRepository,
	storeInfo port.StoreInfo,
	cart port.Cart,
	sender port.MessageSender,
	timeProvider TimeProvider,
	buttonRepository ButtonRepositoryWriter,
) *CartProcessing {
	return &CartProcessing{
		storeID:          store.IDFromInt(storeID),
		repository:       repository,
		storeInfo:        storeInfo,
		cart:             cart,
		sender:           sender,
		timeProvider:     timeProvider,
		buttonRepository: buttonRepository,
	}
}
