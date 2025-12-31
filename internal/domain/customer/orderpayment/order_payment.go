package orderpayment

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
)

type MessageSender interface {
	SendPNG(
		ctx context.Context,
		chatID msginfo.ChatID,
		caption string,
		png []byte,
		rows ...button.ButtonRow,
	) error
	AnswerOrderPayment(
		ctx context.Context,
		paymentID string,
		ok bool,
		errorMsg string,
	) error
	EscapeMarkdown(s string) string
}

type StoreInfo interface {
	GetStoreByID(ctx context.Context, id store.ID) (*store.Store, error)
}

type TimeProvider interface {
	Now() time.Time
}

type VerificationCodeGenerator interface {
	Generate() string
}

type OrderPayment struct {
	storeID       store.ID
	sender        MessageSender
	qrCode        port.QRCodeGenerator
	repository    port.CustomerOrderPaymentRepository
	storeInfo     StoreInfo
	dailyPosition port.DailyPositionGenerator
	codeGenerator VerificationCodeGenerator
	timeProvider  TimeProvider
}

func New(
	storeID int,
	sender MessageSender,
	qrCode port.QRCodeGenerator,
	repository port.CustomerOrderPaymentRepository,
	storeInfo StoreInfo,
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
