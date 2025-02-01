package message

import (
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func OrderNotExists() string {
	return "Order not exists"
}

func OrderTotalPriceIncorrect() string {
	return "Order price incorrect"
}

func OrderStatus(s order.Status) string {
	return fmt.Sprintf("Order in %s status", s.HumanReadable())
}

func OrderStatusChanged(s order.Status) string {
	return fmt.Sprintf("Order status changed to %s", s.HumanReadable())
}

func OrderCanceled() string {
	return "Order canceled"
}

func OrderWithStatusNotExists(s order.Status) string {
	return fmt.Sprintf("Order in %s status is not exists", s.HumanReadable())
}

func OrderInvoice() string {
	return "Order Invoice"
}

func AlreadyHasActiveOrder() string {
	return "You have active order already"
}

func OrderCategoryPage() string {
	return "Select category to view products to order"
}

func CartOrderUnavailable() string {
	return "Order expired or unavailable"
}

func NoProductsForOrder() string {
	return "No products for order"
}

func OrderProductPage() string {
	return "Select product to order"
}

func Cancel() string {
	return "Cancel"
}

func Done() string {
	return "Done"
}

func Confirm() string {
	return "Confirm"
}

func Pay() string {
	return "Pay"
}
