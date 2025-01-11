package button

import (
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

func AddProduct(chatID msginfo.ChatID, productID product.ID) Button {
	return Button{
		ID:        generateID(),
		ChatID:    chatID,
		Operation: OperationCartAddProduct,
		Payload:   []byte(productID.String()),
	}
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
