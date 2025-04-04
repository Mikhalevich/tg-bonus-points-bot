package button

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

type CancelOrderPayload struct {
	OrderID       order.ID
	IsTextMessage bool
}

func CancelOrder(chatID msginfo.ChatID, caption string, id order.ID, isTextMsg bool) (Button, error) {
	return createButton(chatID, caption, OperationOrderCancel,
		CancelOrderPayload{
			OrderID:       id,
			IsTextMessage: isTextMsg,
		},
	)
}

type OrderHistory struct {
	OrderID order.ID
}

func OrderHistoryPrevious(chatID msginfo.ChatID, caption string, beforeID order.ID) (Button, error) {
	return createButton(chatID, caption, OperationOrderHistoryPrevious,
		OrderHistory{
			OrderID: beforeID,
		},
	)
}

func OrderHistoryNext(chatID msginfo.ChatID, caption string, afterID order.ID) (Button, error) {
	return createButton(chatID, caption, OperationOrderHistoryNext,
		OrderHistory{
			OrderID: afterID,
		},
	)
}

func OrderHistoryFirst(chatID msginfo.ChatID, caption string) Button {
	return Button{
		ChatID:    chatID,
		Caption:   caption,
		Operation: OperationOrderHistoryFirst,
	}
}

func OrderHistoryLast(chatID msginfo.ChatID, caption string) Button {
	return Button{
		ChatID:    chatID,
		Caption:   caption,
		Operation: OperationOrderHistoryLast,
	}
}
