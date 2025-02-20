package manager

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
)

type Manager struct {
	customerSender port.MessageSender
	repository     port.ManagerRepository
	timeProvider   port.TimeProvider
}

func New(
	customerSender port.MessageSender,
	repository port.ManagerRepository,
	timeProvider port.TimeProvider,
) *Manager {
	return &Manager{
		customerSender: customerSender,
		repository:     repository,
		timeProvider:   timeProvider,
	}
}
