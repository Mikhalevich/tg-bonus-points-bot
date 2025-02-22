package customerorder

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/store"
)

type CustomerOrder struct {
	storeID          store.ID
	sender           port.MessageSender
	qrCode           port.QRCodeGenerator
	repository       port.CustomerOrderRepository
	storeInfo        port.StoreInfo
	buttonRepository port.ButtonRepository
	dailyPosition    port.DailyPositionGenerator
	codeGenerator    port.VerificationCodeGenerator
	timeProvider     port.TimeProvider
}

func New(
	storeID int,
	sender port.MessageSender,
	qrCode port.QRCodeGenerator,
	repository port.CustomerOrderRepository,
	storeInfo port.StoreInfo,
	buttonRepository port.ButtonRepository,
	dailyPosition port.DailyPositionGenerator,
	codeGenerator port.VerificationCodeGenerator,
	timeProvider port.TimeProvider,
) *CustomerOrder {
	return &CustomerOrder{
		storeID:          store.IDFromInt(storeID),
		sender:           sender,
		qrCode:           qrCode,
		repository:       repository,
		storeInfo:        storeInfo,
		buttonRepository: buttonRepository,
		dailyPosition:    dailyPosition,
		codeGenerator:    codeGenerator,
		timeProvider:     timeProvider,
	}
}
