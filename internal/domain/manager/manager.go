package manager

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
)

type Manager struct {
	customerSender port.MessageSender
	repository     port.ManagerRepository
}

func New(
	customerSender port.MessageSender,
	repository port.ManagerRepository,
) *Manager {
	return &Manager{
		customerSender: customerSender,
		repository:     repository,
	}
}
