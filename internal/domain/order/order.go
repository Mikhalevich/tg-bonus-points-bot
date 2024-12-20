package order

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
)

type Order struct {
	sender port.MessageSender
	qrCode port.QRCodeGenerator
}

func New(
	sender port.MessageSender,
	qrCode port.QRCodeGenerator,
) *Order {
	return &Order{
		sender: sender,
		qrCode: qrCode,
	}
}
