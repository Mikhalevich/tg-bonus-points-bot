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
	ProductID  product.ProductID
	CategoryID product.CategoryID
}

func CartAddProduct(
	chatID msginfo.ChatID,
	caption string,
	cartID cart.ID,
	productID product.ProductID,
	categoryID product.CategoryID,
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
	CategoryID product.CategoryID
}

func CartViewCategoryProducts(
	chatID msginfo.ChatID,
	caption string,
	cartID cart.ID,
	categoryID product.CategoryID,
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
