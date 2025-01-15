package port

import (
	"context"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

type CartItem struct {
	ProductID product.ID
	Count     int
}

type Cart interface {
	StartNewCart(ctx context.Context, chatID msginfo.ChatID) (cart.ID, error)
	GetProducts(ctx context.Context, id cart.ID) ([]CartItem, error)
	AddProduct(ctx context.Context, id cart.ID, productID product.ID) error
	Clear(ctx context.Context, chatID msginfo.ChatID, cartID cart.ID) error
	IsNotFoundError(err error) bool
}
