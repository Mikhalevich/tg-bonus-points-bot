package customercart

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/store"
)

type CustomerCart struct {
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
) *CustomerCart {
	return &CustomerCart{
		storeID:          store.IDFromInt(storeID),
		repository:       repository,
		storeInfo:        storeInfo,
		cart:             cart,
		sender:           sender,
		timeProvider:     timeProvider,
		buttonRepository: buttonRepository,
	}
}
