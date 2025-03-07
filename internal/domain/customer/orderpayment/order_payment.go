package orderpayment

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/store"
)

type OrderPayment struct {
	storeID       store.ID
	sender        port.MessageSender
	qrCode        port.QRCodeGenerator
	repository    port.CustomerOrderPaymentRepository
	storeInfo     port.StoreInfo
	dailyPosition port.DailyPositionGenerator
	codeGenerator port.VerificationCodeGenerator
	timeProvider  port.TimeProvider
}

func New(
	storeID int,
	sender port.MessageSender,
	qrCode port.QRCodeGenerator,
	repository port.CustomerOrderPaymentRepository,
	storeInfo port.StoreInfo,
	dailyPosition port.DailyPositionGenerator,
	codeGenerator port.VerificationCodeGenerator,
	timeProvider port.TimeProvider,
) *OrderPayment {
	return &OrderPayment{
		storeID:       store.IDFromInt(storeID),
		sender:        sender,
		qrCode:        qrCode,
		repository:    repository,
		storeInfo:     storeInfo,
		dailyPosition: dailyPosition,
		codeGenerator: codeGenerator,
		timeProvider:  timeProvider,
	}
}
