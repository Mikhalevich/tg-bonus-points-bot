package button

import (
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func BackToOrder(chatID msginfo.ChatID, id order.ID) Button {
	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationBackToOrder,
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
