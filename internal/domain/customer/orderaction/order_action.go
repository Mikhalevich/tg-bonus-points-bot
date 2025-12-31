package orderaction

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

type MessageSender interface {
	SendMessage(
		ctx context.Context,
		chatID msginfo.ChatID,
		text string,
		textType messageprocessor.MessageTextType,
		rows ...button.ButtonRow,
	) error
	ReplyMessage(
		ctx context.Context,
		chatID msginfo.ChatID,
		replyMessageID msginfo.MessageID,
		text string,
		textType messageprocessor.MessageTextType,
		rows ...button.ButtonRow,
	) error
	EditMessage(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		text string,
		rows ...button.ButtonRow,
	) error
	DeleteMessage(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
	)
	EscapeMarkdown(s string) string
}

type TimeProvider interface {
	Now() time.Time
}

type OrderAction struct {
	sender       MessageSender
	repository   port.CustomerOrderActionRepository
	timeProvider TimeProvider
}

func New(
	sender MessageSender,
	repository port.CustomerOrderActionRepository,
	timeProvider TimeProvider,
) *OrderAction {
	return &OrderAction{
		sender:       sender,
		repository:   repository,
		timeProvider: timeProvider,
	}
}

func (o *OrderAction) sendPlainText(
	ctx context.Context,
	chatID msginfo.ChatID,
	text string,
) {
	if err := o.sender.SendMessage(ctx, chatID, text, messageprocessor.MessageTextTypePlain); err != nil {
		logger.FromContext(ctx).WithError(err).Error("send message")
	}
}

func (o *OrderAction) replyPlainText(
	ctx context.Context,
	chatID msginfo.ChatID,
	replyMessageID msginfo.MessageID,
	text string,
	buttons ...button.ButtonRow,
) {
	if err := o.sender.ReplyMessage(
		ctx,
		chatID,
		replyMessageID,
		text,
		messageprocessor.MessageTextTypePlain,
		buttons...,
	); err != nil {
		logger.FromContext(ctx).WithError(err).Error("reply message")
	}
}

func (o *OrderAction) replyMarkdown(
	ctx context.Context,
	chatID msginfo.ChatID,
	replyMessageID msginfo.MessageID,
	text string,
	buttons ...button.ButtonRow,
) {
	if err := o.sender.ReplyMessage(
		ctx,
		chatID,
		replyMessageID,
		text,
		messageprocessor.MessageTextTypeMarkdown,
		buttons...,
	); err != nil {
		logger.FromContext(ctx).WithError(err).Error("reply message")
	}
}

func (o *OrderAction) editPlainText(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	text string,
	buttons ...button.ButtonRow,
) {
	if err := o.sender.EditMessage(ctx, chatID, messageID, text, buttons...); err != nil {
		logger.FromContext(ctx).WithError(err).Error("edit message")
	}
}
