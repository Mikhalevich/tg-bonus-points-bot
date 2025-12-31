package orderprocessing

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

type CustomerMessageSender interface {
	SendMessage(
		ctx context.Context,
		chatID msginfo.ChatID,
		text string,
		textType messageprocessor.MessageTextType,
		rows ...button.ButtonRow,
	) error
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

func (o *OrderProcessing) sendMarkdown(
	ctx context.Context,
	chatID msginfo.ChatID,
	text string,
	buttons ...button.ButtonRow,
) {
	if err := o.customerSender.SendMessage(
		ctx,
		chatID,
		text,
		messageprocessor.MessageTextTypeMarkdown,
		buttons...,
	); err != nil {
		logger.FromContext(ctx).WithError(err).Error("send message plain")
	}
}
