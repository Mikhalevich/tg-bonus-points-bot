package port

import (
	"context"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

type CartItem struct {
	ProductID product.ID
	Count     int
}

type Cart interface {
	GetProducts(ctx context.Context, chatID msginfo.ChatID) ([]CartItem, error)
	AddProduct(ctx context.Context, chatID msginfo.ChatID, productID product.ID) error
	Clear(ctx context.Context, chatID msginfo.ChatID) error
	IsNotFoundError(err error) bool
}
