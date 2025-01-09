package button

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
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

func ConfirmOrder(chatID msginfo.ChatID, id order.ID) Button {
	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationConfirmOrder,
		Payload:   []byte(id.String()),
	}
}

func (b Button) OrderID() (order.ID, error) {
	id, err := order.IDFromString(string(b.Payload))
	if err != nil {
		return 0, fmt.Errorf("invalid order id: %s", b.Payload)
	}

	return id, nil
}

func ViewCategory(chatID msginfo.ChatID, id product.ID) Button {
	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationViewCategory,
		Payload:   []byte(id.String()),
	}
}

func (b Button) ProductID() (product.ID, error) {
	id, err := product.IDFromString(string(b.Payload))
	if err != nil {
		return 0, fmt.Errorf("invalid id: %w", err)
	}

	return id, nil
}

func Product(chatID msginfo.ChatID, id product.ID) Button {
	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationProduct,
		Payload:   []byte(id.String()),
	}
}

func BackToOrder(chatID msginfo.ChatID, id order.ID) Button {
	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationBackToOrder,
		Payload:   []byte(id.String()),
	}
}
