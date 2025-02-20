package customer

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/store"
)

type Customer struct {
	storeID          store.ID
	sender           port.MessageSender
	qrCode           port.QRCodeGenerator
	repository       port.CustomerRepository
	storeInfo        port.StoreInfo
	cart             port.Cart
	buttonRepository port.ButtonRepository
	dailyPosition    port.DailyPositionGenerator
	codeGenerator    port.VerificationCodeGenerator
	timeProvider     port.TimeProvider
}

func New(
	storeID int,
	sender port.MessageSender,
	qrCode port.QRCodeGenerator,
	repository port.CustomerRepository,
	storeInfo port.StoreInfo,
	cart port.Cart,
	buttonRepository port.ButtonRepository,
	dailyPosition port.DailyPositionGenerator,
	codeGenerator port.VerificationCodeGenerator,
	timeProvider port.TimeProvider,
) *Customer {
	return &Customer{
		storeID:          store.IDFromInt(storeID),
		sender:           sender,
		qrCode:           qrCode,
		repository:       repository,
		storeInfo:        storeInfo,
		cart:             cart,
		buttonRepository: buttonRepository,
		dailyPosition:    dailyPosition,
		codeGenerator:    codeGenerator,
		timeProvider:     timeProvider,
	}
}
