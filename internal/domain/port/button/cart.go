package button

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func CancelCart(chatID msginfo.ChatID) Button {
	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationCartCancel,
	}
}

func ConfirmCart(chatID msginfo.ChatID) Button {
	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationCartConfirm,
	}
}

type AddProductPayload struct {
	ProductID  product.ID
	CategoryID product.ID
}

func AddProduct(chatID msginfo.ChatID, productID, categoryID product.ID) Button {
	var (
		payload = AddProductPayload{
			ProductID:  productID,
			CategoryID: categoryID,
		}

		buf bytes.Buffer
	)

	//nolint:errcheck
	gob.NewEncoder(&buf).Encode(payload)

	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationCartAddProduct,
		Payload:   buf.Bytes(),
	}
}

func (b Button) AddProductPayload() (AddProductPayload, error) {
	var payload AddProductPayload
	if err := gob.NewDecoder(bytes.NewReader(b.Payload)).Decode(&payload); err != nil {
		return AddProductPayload{}, fmt.Errorf("gob decode: %w", err)
	}

	return payload, nil
}

func (b Button) ProductID() (product.ID, error) {
	id, err := product.IDFromString(string(b.Payload))
	if err != nil {
		return 0, fmt.Errorf("invalid payload: %w", err)
	}

	return id, nil
}

func ViewCategoryProducts(chatID msginfo.ChatID, categoryID product.ID) Button {
	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationCartViewCategoryProducts,
		Payload:   []byte(categoryID.String()),
	}
}

func (b Button) CategoryID() (product.ID, error) {
	id, err := product.IDFromString(string(b.Payload))
	if err != nil {
		return 0, fmt.Errorf("invalid payload: %w", err)
	}

	return id, nil
}

func ViewCategories(chatID msginfo.ChatID) Button {
	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationCartViewCategories,
	}
}
