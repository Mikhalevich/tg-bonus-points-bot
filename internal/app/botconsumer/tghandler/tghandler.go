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

type CartProcessor interface {
	Create(ctx context.Context, info msginfo.Info) error
	ViewCategoryProducts(ctx context.Context, info msginfo.Info, cartID cart.ID, categoryID product.CategoryID,
		currencyID currency.ID) error
	ViewCategories(ctx context.Context, info msginfo.Info, cartID cart.ID, currencyID currency.ID) error
	AddProduct(ctx context.Context, info msginfo.Info, cartID cart.ID, categoryID product.CategoryID,
		productID product.ProductID, currencyID currency.ID) error
	Cancel(ctx context.Context, info msginfo.Info, cartID cart.ID) error
	Confirm(ctx context.Context, info msginfo.Info, cartID cart.ID, currencyID currency.ID) error
}

type OrderActionProcessor interface {
	GetActiveOrder(ctx context.Context, info msginfo.Info) error
	Cancel(ctx context.Context, chatID msginfo.ChatID, messageID msginfo.MessageID,
		orderID order.ID, isTextMsg bool) error
	QueueSize(ctx context.Context, info msginfo.Info) error
}

type OrderHistoryProcessor interface {
	Show(ctx context.Context, chatID msginfo.ChatID) error
	First(ctx context.Context, info msginfo.Info) error
	Last(ctx context.Context, info msginfo.Info) error
	Previous(ctx context.Context, info msginfo.Info, beforeOrderID order.ID) error
	Next(ctx context.Context, info msginfo.Info, afterOrderID order.ID) error
}

type OrderPaymentProcessor interface {
	PaymentInProgress(ctx context.Context, paymentID string, orderID order.ID,
		currency string, totalAmount int) error
	PaymentConfirmed(ctx context.Context, chatID msginfo.ChatID, orderID order.ID,
		currency string, totalAmount int) error
}

type ButtonProvider interface {
	GetButton(ctx context.Context, id button.ID) (*button.Button, error)
}

type cbHandler func(ctx context.Context, info msginfo.Info, btn button.Button) error

type TGHandler struct {
	cartProcessor    CartProcessor
	actionProcessor  OrderActionProcessor
	historyProcessor OrderHistoryProcessor
	paymentProcessor OrderPaymentProcessor
	buttonProvider   ButtonProvider
	cbHandlers       map[button.Operation]cbHandler
}

func New(
	cartProcessor CartProcessor,
	actionProcessor OrderActionProcessor,
	historyProcessor OrderHistoryProcessor,
	paymentProcessor OrderPaymentProcessor,
	buttonProvider ButtonProvider,
) *TGHandler {
	handler := &TGHandler{
		cartProcessor:    cartProcessor,
		actionProcessor:  actionProcessor,
		historyProcessor: historyProcessor,
		paymentProcessor: paymentProcessor,
		buttonProvider:   buttonProvider,
	}

	handler.initCBHandlers()

	return handler
}

func (t *TGHandler) initCBHandlers() {
	t.cbHandlers = map[button.Operation]cbHandler{
		button.OperationOrderCancel:              t.cancelOrder,
		button.OperationCartCancel:               t.cancelCart,
		button.OperationCartConfirm:              t.confirmCart,
		button.OperationCartViewCategories:       t.viewCategories,
		button.OperationCartViewCategoryProducts: t.viewCategoryProducts,
		button.OperationCartAddProduct:           t.addProduct,
		button.OperationOrderHistoryPrevious:     t.historyPrevious,
		button.OperationOrderHistoryNext:         t.historyNext,
		button.OperationOrderHistoryFirst:        t.historyFirst,
		button.OperationOrderHistoryLast:         t.historyLast,
	}
}
