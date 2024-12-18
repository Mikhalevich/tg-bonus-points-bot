package order

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
)

type Order struct {
	sender port.MessageSender
}

func New(sender port.MessageSender) *Order {
	return &Order{
		sender: sender,
	}
}
