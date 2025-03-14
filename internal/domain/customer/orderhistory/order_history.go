package orderhistory

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
)

type OrderHistory struct {
	repository       port.CustomerOrderHistoryRepository
	sender           port.MessageSender
	buttonRepository port.ButtonRepositoryWriter
	pageSize         int
}

func New(
	repository port.CustomerOrderHistoryRepository,
	sender port.MessageSender,
	buttonRepository port.ButtonRepositoryWriter,
	pageSize int,
) *OrderHistory {
	return &OrderHistory{
		repository:       repository,
		sender:           sender,
		buttonRepository: buttonRepository,
		pageSize:         pageSize,
	}
}
