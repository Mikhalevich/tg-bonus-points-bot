package button

import (
	"fmt"

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

func CancelOrder(chatID msginfo.ChatID, id order.ID) Button {
	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationCancelOrder,
		Payload:   []byte(id.String()),
	}
}

func (b Button) OrderID() (order.ID, error) {
	id, err := order.IDFromString(string(b.Payload))
	if err != nil {
		return 0, fmt.Errorf("invaid order id: %w", err)
	}

	return id, nil
}

func CancelCart(chatID msginfo.ChatID) Button {
	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationCancelCart,
	}
}

func ConfirmCart(chatID msginfo.ChatID) Button {
	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationConfirmCart,
	}
}
