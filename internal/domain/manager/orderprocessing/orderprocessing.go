package orderprocessing

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

type CustomerMessageSender interface {
	SendTextMarkdown(
		ctx context.Context,
		chatID msginfo.ChatID,
		text string,
		buttons ...button.InlineKeyboardButtonRow,
	)
	EscapeMarkdown(s string) string
}

type Repository interface {
	UpdateOrderStatusForMinID(
		ctx context.Context,
		operationTime time.Time,
		newStatus, prevStatus order.Status,
	) (*order.Order, error)
	UpdateOrderStatus(
		ctx context.Context,
		id order.ID,
		operationTime time.Time,
		newStatus order.Status,
		prevStatuses ...order.Status,
	) (*order.Order, error)
	IsNotFoundError(err error) bool
	IsNotUpdatedError(err error) bool
}

type TimeProvider interface {
	Now() time.Time
}

type OrderProcessing struct {
	customerSender CustomerMessageSender
	repository     Repository
	timeProvider   TimeProvider
}

func New(
	customerSender CustomerMessageSender,
	repository Repository,
	timeProvider TimeProvider,
) *OrderProcessing {
	return &OrderProcessing{
		customerSender: customerSender,
		repository:     repository,
		timeProvider:   timeProvider,
	}
}
