package port

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

type ManagerRepository interface {
	UpdateOrderStatusForMinID(
		ctx context.Context,
		operationTime time.Time,
		prevStatus, newStatus order.Status,
	) (*order.Order, error)
	IsNotFoundError(err error) bool
}
