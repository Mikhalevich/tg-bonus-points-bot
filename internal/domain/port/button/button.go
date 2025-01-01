package button

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

type Operation string

const (
	OperationCancelOrder Operation = "CancelOrder"
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
	Data      any
}

func (b Button) OrderID() (order.ID, error) {
	id, ok := b.Data.(order.ID)
	if !ok {
		return 0, fmt.Errorf("invalid order id: %v", b.Data)
	}

	return id, nil
}

func CancelOrder(chatID msginfo.ChatID, orderID order.ID) Button {
	return Button{
		ID:        IDFromString(uuid.NewString()),
		ChatID:    chatID,
		Operation: OperationCancelOrder,
		Data:      orderID,
	}
}
