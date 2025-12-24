package orderaction

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
)

type ButtonRepositoryWriter interface {
	SetButton(ctx context.Context, b button.Button) (button.InlineKeyboardButton, error)
}

type OrderAction struct {
	sender           port.MessageSender
	repository       port.CustomerOrderActionRepository
	buttonRepository ButtonRepositoryWriter
	timeProvider     port.TimeProvider
}

func New(
	sender port.MessageSender,
	repository port.CustomerOrderActionRepository,
	buttonRepository ButtonRepositoryWriter,
	timeProvider port.TimeProvider,
) *OrderAction {
	return &OrderAction{
		sender:           sender,
		repository:       repository,
		buttonRepository: buttonRepository,
		timeProvider:     timeProvider,
	}
}
