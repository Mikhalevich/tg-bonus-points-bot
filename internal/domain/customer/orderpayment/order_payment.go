package orderpayment

import (
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
)

type TimeProvider interface {
	Now() time.Time
}

type VerificationCodeGenerator interface {
	Generate() string
}

type OrderPayment struct {
	storeID       store.ID
	sender        port.MessageSender
	qrCode        port.QRCodeGenerator
	repository    port.CustomerOrderPaymentRepository
	storeInfo     port.StoreInfo
	dailyPosition port.DailyPositionGenerator
	codeGenerator VerificationCodeGenerator
	timeProvider  TimeProvider
}

func New(
	storeID int,
	sender port.MessageSender,
	qrCode port.QRCodeGenerator,
	repository port.CustomerOrderPaymentRepository,
	storeInfo port.StoreInfo,
	dailyPosition port.DailyPositionGenerator,
	codeGenerator VerificationCodeGenerator,
	timeProvider TimeProvider,
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
