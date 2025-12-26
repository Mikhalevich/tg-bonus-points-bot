package orderaction

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
)

type ButtonRepositoryWriter interface {
	SetButton(ctx context.Context, b button.Button) (button.InlineKeyboardButton, error)
}

type TimeProvider interface {
	Now() time.Time
}

type OrderAction struct {
	sender           port.MessageSender
	repository       port.CustomerOrderActionRepository
	buttonRepository ButtonRepositoryWriter
	timeProvider     TimeProvider
}

func New(
	sender port.MessageSender,
	repository port.CustomerOrderActionRepository,
	buttonRepository ButtonRepositoryWriter,
	timeProvider TimeProvider,
) *OrderAction {
	return &OrderAction{
		sender:           sender,
		repository:       repository,
		buttonRepository: buttonRepository,
		timeProvider:     timeProvider,
	}
}
