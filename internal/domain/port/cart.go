package port

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

type Cart interface {
	StartNewCart(ctx context.Context, chatID msginfo.ChatID) (cart.ID, error)
	GetProducts(ctx context.Context, id cart.ID) ([]cart.CartProduct, error)
	AddProduct(ctx context.Context, id cart.ID, p cart.CartProduct) error
	Clear(ctx context.Context, chatID msginfo.ChatID, cartID cart.ID) error
	IsNotFoundError(err error) bool
}
