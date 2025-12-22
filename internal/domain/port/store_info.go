package port

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
)

type StoreInfo interface {
	GetStoreByID(ctx context.Context, id store.ID) (*store.Store, error)
}
