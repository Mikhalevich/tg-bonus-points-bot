package orderhistory

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
)

type OrderHistory struct {
	repository port.CustomerOrderHistoryRepository
	sender     port.MessageSender
}

func New(
	repository port.CustomerOrderHistoryRepository,
	sender port.MessageSender,
) *OrderHistory {
	return &OrderHistory{
		repository: repository,
		sender:     sender,
	}
}
