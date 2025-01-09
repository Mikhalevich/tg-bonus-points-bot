package button

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

type ViewCategoryPayload struct {
	OrderID    order.ID
	CategoryID product.ID
}

func ViewCategory(chatID msginfo.ChatID, orderID order.ID, categoryID product.ID) Button {
	var (
		buf     bytes.Buffer
		payload = ViewCategoryPayload{
			OrderID:    orderID,
			CategoryID: categoryID,
		}
	)

	//nolint:errcheck
	gob.NewEncoder(&buf).Encode(payload)

	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationViewCategory,
		Payload:   buf.Bytes(),
	}
}

func (b Button) ViewCategoryPayload() (ViewCategoryPayload, error) {
	var payload ViewCategoryPayload
	if err := gob.NewDecoder(bytes.NewReader(b.Payload)).Decode(&payload); err != nil {
		return ViewCategoryPayload{}, fmt.Errorf("decode payload: %w", err)
	}

	return payload, nil
}
