package button

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

func BackToOrder(chatID msginfo.ChatID) Button {
	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationBackToOrder,
	}
}
