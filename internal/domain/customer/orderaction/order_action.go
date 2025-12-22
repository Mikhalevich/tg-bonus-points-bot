package orderaction

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port"
)

type OrderAction struct {
	sender           port.MessageSender
	repository       port.CustomerOrderActionRepository
	buttonRepository port.ButtonRepositoryWriter
	timeProvider     port.TimeProvider
}

func New(
	sender port.MessageSender,
	repository port.CustomerOrderActionRepository,
	buttonRepository port.ButtonRepositoryWriter,
	timeProvider port.TimeProvider,
) *OrderAction {
	return &OrderAction{
		sender:           sender,
		repository:       repository,
		buttonRepository: buttonRepository,
		timeProvider:     timeProvider,
	}
}
