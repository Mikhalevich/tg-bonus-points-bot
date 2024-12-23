package order

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
)

type Order struct {
	sender     port.MessageSender
	qrCode     port.QRCodeGenerator
	repository port.OrderRepository
}

func New(
	sender port.MessageSender,
	qrCode port.QRCodeGenerator,
	repository port.OrderRepository,
) *Order {
	return &Order{
		sender:     sender,
		qrCode:     qrCode,
		repository: repository,
	}
}
