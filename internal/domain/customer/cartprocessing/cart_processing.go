package cartprocessing

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
)

type CartProcessing struct {
	storeID          store.ID
	repository       port.CustomerCartRepository
	storeInfo        port.StoreInfo
	cart             port.Cart
	sender           port.MessageSender
	timeProvider     port.TimeProvider
	buttonRepository port.ButtonRepositoryWriter
}

func New(
	storeID int,
	repository port.CustomerCartRepository,
	storeInfo port.StoreInfo,
	cart port.Cart,
	sender port.MessageSender,
	timeProvider port.TimeProvider,
	buttonRepository port.ButtonRepositoryWriter,
) *CartProcessing {
	return &CartProcessing{
		storeID:          store.IDFromInt(storeID),
		repository:       repository,
		storeInfo:        storeInfo,
		cart:             cart,
		sender:           sender,
		timeProvider:     timeProvider,
		buttonRepository: buttonRepository,
	}
}
