package tghandler

import (
	"context"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

//nolint:interfacebloat
type OrderProcessor interface {
	CartCreate(ctx context.Context, info msginfo.Info) error
	CartViewCategoryProducts(ctx context.Context, info msginfo.Info, cartID cart.ID, categoryID product.CategoryID,
		currencyID currency.ID) error
	CartViewCategories(ctx context.Context, info msginfo.Info, cartID cart.ID, currencyID currency.ID) error
	CartAddProduct(ctx context.Context, info msginfo.Info, cartID cart.ID, categoryID product.CategoryID,
		productID product.ProductID, currencyID currency.ID) error
	CartCancel(ctx context.Context, info msginfo.Info, cartID cart.ID) error
	CartConfirm(ctx context.Context, info msginfo.Info, cartID cart.ID, currencyID currency.ID) error

	GetButton(ctx context.Context, id button.ID) (*button.Button, error)

	GetActiveOrder(ctx context.Context, info msginfo.Info) error
	OrderCancel(ctx context.Context, chatID msginfo.ChatID, messageID msginfo.MessageID,
		orderID order.ID, isTextMsg bool) error
	OrderPaymentInProgress(ctx context.Context, paymentID string, orderID order.ID,
		currency string, totalAmount int) error
	OrderPaymentConfirmed(ctx context.Context, chatID msginfo.ChatID, orderID order.ID,
		currency string, totalAmount int) error
	OrderQueueSize(ctx context.Context, info msginfo.Info) error
}

type cbHandler func(ctx context.Context, info msginfo.Info, btn button.Button) error

type TGHandler struct {
	orderProcessor OrderProcessor
	cbHandlers     map[button.Operation]cbHandler
}

func New(orderProcessor OrderProcessor) *TGHandler {
	h := &TGHandler{
		orderProcessor: orderProcessor,
	}

	h.initCBHandlers()

	return h
}

func (t *TGHandler) initCBHandlers() {
	t.cbHandlers = map[button.Operation]cbHandler{
		button.OperationOrderCancel:              t.cancelOrder,
		button.OperationCartCancel:               t.cancelCart,
		button.OperationCartConfirm:              t.confirmCart,
		button.OperationCartViewCategories:       t.viewCategories,
		button.OperationCartViewCategoryProducts: t.viewCategoryProducts,
		button.OperationCartAddProduct:           t.addProduct,
	}
}
