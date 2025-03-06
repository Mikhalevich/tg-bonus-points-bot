package port

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

type UpdateOrderData struct {
	Status              order.Status
	StatusOperationTime time.Time
	VerificationCode    string
	DailyPosition       int
}

//nolint:interfacebloat
type CustomerOrderRepository interface {
	GetOrderByChatIDAndStatus(ctx context.Context, id msginfo.ChatID, statuses ...order.Status) (*order.Order, error)
	GetOrderByID(ctx context.Context, id order.ID) (*order.Order, error)
	UpdateOrderByChatAndID(
		ctx context.Context,
		orderID order.ID,
		chatID msginfo.ChatID,
		data UpdateOrderData,
		prevStatuses ...order.Status,
	) (*order.Order, error)
	UpdateOrderStatusByChatAndID(
		ctx context.Context,
		orderID order.ID,
		chatID msginfo.ChatID,
		operationTime time.Time,
		newStatus order.Status,
		prevStatuses ...order.Status,
	) (*order.Order, error)
	UpdateOrderStatus(
		ctx context.Context,
		id order.ID,
		operationTime time.Time,
		newStatus order.Status,
		prevStatuses ...order.Status,
	) (*order.Order, error)
	GetOrdersCountByStatus(ctx context.Context, statuses ...order.Status) (int, error)
	GetOrderPositionByStatus(ctx context.Context, id order.ID, statuses ...order.Status) (int, error)
	GetProductsByIDs(
		ctx context.Context,
		ids []product.ProductID,
		currencyID currency.ID,
	) (map[product.ProductID]product.Product, error)
	HistoryOrders(ctx context.Context, chatID msginfo.ChatID, size int) ([]order.ShortOrder, error)
	IsNotFoundError(err error) bool
	IsNotUpdatedError(err error) bool
}
