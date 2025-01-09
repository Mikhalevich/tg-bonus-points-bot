package button

import (
	"github.com/google/uuid"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

type ID string

func (id ID) String() string {
	return string(id)
}

func IDFromString(s string) ID {
	return ID(s)
}

type Button struct {
	ID        ID
	ChatID    msginfo.ChatID
	Operation Operation
	Payload   []byte
}

func generateID() ID {
	return IDFromString(uuid.NewString())
}

func CancelOrderSendMsg(chatID msginfo.ChatID, id order.ID) Button {
	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationCancelOrderSendMessage,
		Payload:   []byte(id.String()),
	}
}

func CancelOrderEditMsg(chatID msginfo.ChatID, id order.ID) Button {
	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationCancelOrderEditMessage,
		Payload:   []byte(id.String()),
	}
}

func ConfirmOrder(chatID msginfo.ChatID, id order.ID) Button {
	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationConfirmOrder,
		Payload:   []byte(id.String()),
	}
}
