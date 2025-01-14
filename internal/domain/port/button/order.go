package button

import (
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func CancelOrder(chatID msginfo.ChatID, caption string, id order.ID) Button {
	return Button{
		ChatID:    chatID,
		Caption:   caption,
		Operation: OperationOrderCancel,
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
