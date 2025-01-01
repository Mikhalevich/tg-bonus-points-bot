package customer

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
)

type Customer struct {
	sender           port.MessageSender
	qrCode           port.QRCodeGenerator
	repository       port.CustomerRepository
	buttonRepository port.ButtonRepository
}

func New(
	sender port.MessageSender,
	qrCode port.QRCodeGenerator,
	repository port.CustomerRepository,
	buttonRepository port.ButtonRepository,
) *Customer {
	return &Customer{
		sender:           sender,
		qrCode:           qrCode,
		repository:       repository,
		buttonRepository: buttonRepository,
	}
}
