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

type OrderHistoryByID struct {
	OrderID order.ID
}

func OrderHistoryByIDPrevious(chatID msginfo.ChatID, caption string, beforeID order.ID) (Button, error) {
	return createButton(chatID, caption, OperationOrderHistoryByIDPrevious,
		OrderHistoryByID{
			OrderID: beforeID,
		},
	)
}

func OrderHistoryByIDNext(chatID msginfo.ChatID, caption string, afterID order.ID) (Button, error) {
	return createButton(chatID, caption, OperationOrderHistoryByIDNext,
		OrderHistoryByID{
			OrderID: afterID,
		},
	)
}

func OrderHistoryByIDFirst(chatID msginfo.ChatID, caption string) Button {
	return Button{
		ChatID:    chatID,
		Caption:   caption,
		Operation: OperationOrderHistoryByIDFirst,
	}
}

func OrderHistoryByIDLast(chatID msginfo.ChatID, caption string) Button {
	return Button{
		ChatID:    chatID,
		Caption:   caption,
		Operation: OperationOrderHistoryByIDLast,
	}
}

type OrderHistoryByPagePayload struct {
	Page int
}

func OrderHistoryByPage(chatID msginfo.ChatID, caption string, page int) (Button, error) {
	return createButton(chatID, caption, OperationOrderHistoryByPage,
		OrderHistoryByPagePayload{
			Page: page,
		},
	)
}

func OrderHistoryByPageFirst(chatID msginfo.ChatID, caption string) Button {
	return Button{
		ChatID:    chatID,
		Caption:   caption,
		Operation: OperationOrderHistoryByPageFirst,
	}
}

func OrderHistoryByPageLast(chatID msginfo.ChatID, caption string) Button {
	return Button{
		ChatID:    chatID,
		Caption:   caption,
		Operation: OperationOrderHistoryByPageLast,
	}
}
