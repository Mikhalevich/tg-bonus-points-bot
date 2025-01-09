package message

import (
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func OrderNotExists() string {
	return "Order not exists"
}

func OrderStatus(s order.Status) string {
	return fmt.Sprintf("Order in %s status", s.HumanReadable())
}

func OrderStatusChanged(s order.Status) string {
	return fmt.Sprintf("Order status changed to %s", s.HumanReadable())
}

func OrderWithStatusNotExists(s order.Status) string {
	return fmt.Sprintf("Order in %s status is not exists", s.HumanReadable())
}

func OrderCategoryPage() string {
	return "Choose category to view products to order"
}

func Cancel() string {
	return "Cancel"
}
