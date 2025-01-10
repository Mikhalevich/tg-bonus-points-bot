package button

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

type ProductPayload struct {
	ProductID product.ID
}

func Product(chatID msginfo.ChatID, productID product.ID) Button {
	var (
		buf     bytes.Buffer
		payload = ViewCategoryPayload{
			CategoryID: productID,
		}
	)

	//nolint:errcheck
	gob.NewEncoder(&buf).Encode(payload)

	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationProduct,
		Payload:   buf.Bytes(),
	}
}

func (b Button) ProductPayload() (ProductPayload, error) {
	var payload ProductPayload
	if err := gob.NewDecoder(bytes.NewReader(b.Payload)).Decode(&payload); err != nil {
		return ProductPayload{}, fmt.Errorf("decode payload: %w", err)
	}

	return payload, nil
}
