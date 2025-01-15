package button

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

type CartCancelPayload struct {
	CartID cart.ID
}

func CartCancel(chatID msginfo.ChatID, caption string, cartID cart.ID) (Button, error) {
	return createButton(chatID, caption, OperationCartCancel,
		CartCancelPayload{
			CartID: cartID,
		},
	)
}

type CartConfirmPayload struct {
	CartID cart.ID
}

func CartConfirm(chatID msginfo.ChatID, caption string, cartID cart.ID) (Button, error) {
	return createButton(chatID, caption, OperationCartConfirm,
		CartConfirmPayload{
			CartID: cartID,
		},
	)
}

type CartAddProductPayload struct {
	CartID     cart.ID
	ProductID  product.ID
	CategoryID product.ID
}

func CartAddProduct(
	chatID msginfo.ChatID,
	caption string,
	cartID cart.ID,
	productID, categoryID product.ID,
) (Button, error) {
	return createButton(chatID, caption, OperationCartAddProduct,
		CartAddProductPayload{
			CartID:     cartID,
			ProductID:  productID,
			CategoryID: categoryID,
		},
	)
}

type CartViewCategoryProductsPayload struct {
	CartID     cart.ID
	CategoryID product.ID
}

func CartViewCategoryProducts(
	chatID msginfo.ChatID,
	caption string,
	cartID cart.ID,
	categoryID product.ID,
) (Button, error) {
	return createButton(chatID, caption, OperationCartViewCategoryProducts,
		CartViewCategoryProductsPayload{
			CartID:     cartID,
			CategoryID: categoryID,
		},
	)
}

type CartViewCategoriesPayload struct {
	CartID cart.ID
}

func CartViewCategories(chatID msginfo.ChatID, caption string, cartID cart.ID) (Button, error) {
	return createButton(chatID, caption, OperationCartViewCategories,
		CartViewCategoriesPayload{
			CartID: cartID,
		},
	)
}
