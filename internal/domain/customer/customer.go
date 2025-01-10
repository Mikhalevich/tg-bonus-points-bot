package customer

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
)

type Customer struct {
	sender           port.MessageSender
	qrCode           port.QRCodeGenerator
	repository       port.CustomerRepository
	cart             port.Cart
	buttonRepository port.ButtonRepository
}

func New(
	sender port.MessageSender,
	qrCode port.QRCodeGenerator,
	repository port.CustomerRepository,
	cart port.Cart,
	buttonRepository port.ButtonRepository,
) *Customer {
	return &Customer{
		sender:           sender,
		qrCode:           qrCode,
		repository:       repository,
		cart:             cart,
		buttonRepository: buttonRepository,
	}
}
