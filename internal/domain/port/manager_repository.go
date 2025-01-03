package port

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

type ManagerRepository interface {
	UpdateOrderStatusForMinID(
		ctx context.Context,
		operationTime time.Time,
		prevStatus, newStatus order.Status,
	) (*order.Order, error)
	UpdateOrderStatus(
		ctx context.Context,
		id order.ID,
		operationTime time.Time,
		newStatus order.Status,
		prevStatuses ...order.Status,
	) (*order.Order, error)
	GetCategoryProducts(ctx context.Context) ([]product.Category, error)
	IsNotFoundError(err error) bool
	IsNotUpdatedError(err error) bool
}
