package button

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func CancelCart(chatID msginfo.ChatID, caption string) Button {
	return Button{
		ChatID:    chatID,
		Caption:   caption,
		Operation: OperationCartCancel,
	}
}

type CartConfirmPayload struct {
	CartID cart.ID
}

func (b Button) CartConfirmPayload() (CartConfirmPayload, error) {
	payload, err := gobDecodePayload[CartConfirmPayload](b.Payload)
	if err != nil {
		return CartConfirmPayload{}, fmt.Errorf("decode payload: %w", err)
	}

	return payload, nil
}

func ConfirmCart(chatID msginfo.ChatID, caption string, cartID cart.ID) (Button, error) {
	payload, err := gobEncodePayload(CartConfirmPayload{
		CartID: cartID,
	})

	if err != nil {
		return Button{}, fmt.Errorf("encode payload: %w", err)
	}

	return Button{
		ChatID:    chatID,
		Caption:   caption,
		Operation: OperationCartConfirm,
		Payload:   payload,
	}, nil
}

type AddProductPayload struct {
	CartID     cart.ID
	ProductID  product.ID
	CategoryID product.ID
}

func AddProduct(
	chatID msginfo.ChatID,
	caption string,
	cartID cart.ID,
	productID, categoryID product.ID,
) (Button, error) {
	payload, err := gobEncodePayload(AddProductPayload{
		CartID:     cartID,
		ProductID:  productID,
		CategoryID: categoryID,
	})

	if err != nil {
		return Button{}, fmt.Errorf("encode payload: %w", err)
	}

	return Button{
		ChatID:    chatID,
		Caption:   caption,
		Operation: OperationCartAddProduct,
		Payload:   payload,
	}, nil
}

func gobEncodePayload(p any) ([]byte, error) {
	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(p); err != nil {
		return nil, fmt.Errorf("gob encode: %w", err)
	}

	return buf.Bytes(), nil
}

func gobDecodePayload[Payload any](b []byte) (Payload, error) {
	var payload Payload
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&payload); err != nil {
		return payload, fmt.Errorf("gob decode: %w", err)
	}

	return payload, nil
}

func (b Button) AddProductPayload() (AddProductPayload, error) {
	payload, err := gobDecodePayload[AddProductPayload](b.Payload)
	if err != nil {
		return AddProductPayload{}, fmt.Errorf("decode payload: %w", err)
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

type CartCategoryProductsPayload struct {
	CartID     cart.ID
	CategoryID product.ID
}

func CartViewCategoryProducts(
	chatID msginfo.ChatID,
	caption string,
	cartID cart.ID,
	categoryID product.ID,
) (Button, error) {
	payload, err := gobEncodePayload(CartCategoryProductsPayload{
		CartID:     cartID,
		CategoryID: categoryID,
	})

	if err != nil {
		return Button{}, fmt.Errorf("encode payload: %w", err)
	}

	return Button{
		ChatID:    chatID,
		Caption:   caption,
		Operation: OperationCartViewCategoryProducts,
		Payload:   payload,
	}, nil
}

func (b Button) CartCategoryProductsPayload() (CartCategoryProductsPayload, error) {
	payload, err := gobDecodePayload[CartCategoryProductsPayload](b.Payload)
	if err != nil {
		return CartCategoryProductsPayload{}, fmt.Errorf("decode payload: %w", err)
	}

	return payload, nil
}

type CartViewCategoriesPayload struct {
	CartID cart.ID
}

func CartViewCategories(chatID msginfo.ChatID, caption string, cartID cart.ID) (Button, error) {
	payload, err := gobEncodePayload(CartViewCategoriesPayload{
		CartID: cartID,
	})

	if err != nil {
		return Button{}, fmt.Errorf("encode payload: %w", err)
	}

	return Button{
		ChatID:    chatID,
		Caption:   caption,
		Operation: OperationCartViewCategories,
		Payload:   payload,
	}, nil
}

func (b Button) CartViewCategoriesPayload() (CartViewCategoriesPayload, error) {
	payload, err := gobDecodePayload[CartViewCategoriesPayload](b.Payload)
	if err != nil {
		return CartViewCategoriesPayload{}, fmt.Errorf("decode payload: %w", err)
	}

	return payload, nil
}
