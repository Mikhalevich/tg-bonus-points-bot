package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/msgconsumer/internal/app/event"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

type Consumer interface {
	Consume(
		ctx context.Context,
		processFn func(ctx context.Context, payload []byte) error,
	) error
}

type MessageSender interface {
	SendMessage(
		ctx context.Context,
		chatID msginfo.ChatID,
		text string,
		textType messageprocessor.MessageTextType,
		rows ...button.ButtonRow,
	) error
}

type App struct {
	consumer Consumer
	sender   MessageSender
}

func New(
	consumer Consumer,
	sender MessageSender,
) *App {
	return &App{
		consumer: consumer,
		sender:   sender,
	}
}

func (a *App) Start(ctx context.Context) error {
	if err := a.consumer.Consume(ctx, func(ctx context.Context, payload []byte) error {
		var msg event.OutboxMessage
		if err := json.Unmarshal(payload, &msg); err != nil {
			return fmt.Errorf("unmarshal message: %w", err)
		}

		switch msg.MessageType {
		case event.OutboxMessageTypePlain:
			return a.sendTextMessage(ctx, &msg, messageprocessor.MessageTextTypePlain)

		case event.OutboxMessageTypeMarkdown:
			return a.sendTextMessage(ctx, &msg, messageprocessor.MessageTextTypeMarkdown)
		}

		return fmt.Errorf("unknown message type: %s", msg.MessageType)
	}); err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	return nil
}

func (a *App) sendTextMessage(
	ctx context.Context,
	msg *event.OutboxMessage,
	msgType messageprocessor.MessageTextType,
) error {
	if err := a.sender.SendMessage(
		ctx,
		msginfo.ChatID(msg.ChatID),
		msg.MessageText,
		msgType,
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
