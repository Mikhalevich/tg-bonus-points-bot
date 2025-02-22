package port

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/currency"
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
	CurrencyID          currency.ID
}

type CustomerCartRepository interface {
	CreateOrder(ctx context.Context, coi CreateOrderInput) (*order.Order, error)
	GetCategories(ctx context.Context) ([]product.Category, error)
	GetProductsByCategoryID(
		ctx context.Context,
		categoryID product.CategoryID,
		currencyID currency.ID,
	) ([]product.Product, error)
	GetProductsByIDs(
		ctx context.Context,
		ids []product.ProductID,
		currencyID currency.ID,
	) (map[product.ProductID]product.Product, error)
	GetCurrencyByID(ctx context.Context, id currency.ID) (*currency.Currency, error)
	IsAlreadyExistsError(err error) bool
}
