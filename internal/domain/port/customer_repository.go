package port

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

type CreateOrderInput struct {
	ChatID              msginfo.ChatID
	Status              order.Status
	StatusOperationTime time.Time
	VerificationCode    string
	Products            []order.OrderedProduct
}

type CustomerRepository interface {
	CreateOrder(ctx context.Context, coi CreateOrderInput) (order.ID, error)
	GetOrderByChatIDAndStatus(ctx context.Context, id msginfo.ChatID, statuses ...order.Status) (*order.Order, error)
	GetOrderByID(ctx context.Context, id order.ID) (*order.Order, error)
	UpdateOrderStatus(
		ctx context.Context,
		id order.ID,
		operationTime time.Time,
		newStatus order.Status,
		prevStatuses ...order.Status,
	) (*order.Order, error)
	GetCategories(ctx context.Context) ([]product.Category, error)
	GetProductsByCategoryID(ctx context.Context, id product.CategoryID) ([]product.Product, error)
	GetProductsByIDs(ctx context.Context, ids []product.ProductID) (map[product.ProductID]product.Product, error)
	IsNotFoundError(err error) bool
	IsNotUpdatedError(err error) bool
	IsAlreadyExistsError(err error) bool
}
