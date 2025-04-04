package tgbot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/tracing"
)

type Payment struct {
	IsSuccessful   bool
	IsCheckout     bool
	ID             string
	InvoicePayload string
	Currency       string
	TotalAmount    int
}

type BotMessage struct {
	MessageID int
	ChatID    int64
	// for text message
	Text string
	// for callback query
	Data    string
	Payment Payment
}

type MessageSender interface {
	SendMessage(ctx context.Context, chatID int64, msg string)
}

type Handler func(ctx context.Context, msg BotMessage, sender MessageSender) error

func (t *TGBot) AddMenuCommand(command string, description string, handler Handler) {
	t.addCommand(command, description, handler)
}

func (t *TGBot) AddTextCommand(command string, handler Handler) {
	t.addCommand(command, "", handler)
}

func (t *TGBot) addCommand(command string, description string, handler Handler) {
	if description != "" {
		t.commands = append(t.commands, models.BotCommand{
			Command:     command,
			Description: description,
		})
	}

	t.bot.RegisterHandler(
		bot.HandlerTypeMessageText,
		command,
		bot.MatchTypeExact,
		t.wrapHandler(command, handler),
	)
}

func (t *TGBot) AddDefaultHandler(h Handler) {
	h = t.applyMiddleware(h)
	t.defaultHandlerFn = h
}

func (t *TGBot) AddDefaultTextHandler(h Handler) {
	t.bot.RegisterHandler(
		bot.HandlerTypeMessageText,
		"",
		bot.MatchTypePrefix,
		t.wrapHandler("default_text_handler", h),
	)
}

func (t *TGBot) AddDefaultCallbackQueryHander(h Handler) {
	t.bot.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"",
		bot.MatchTypePrefix,
		t.wrapHandler("default_callback_query", h),
	)
}

func (t *TGBot) wrapHandler(pattern string, handler Handler) bot.HandlerFunc {
	handler = t.applyMiddleware(handler)

	return func(ctx context.Context, botAPI *bot.Bot, update *models.Update) {
		ctx, span := tracing.StartSpanName(ctx, pattern)
		defer span.End()

		var (
			msg    = makeMsgFromUpdate(update)
			log    = t.logger.WithContext(ctx).WithField("endpoint", pattern)
			ctxLog = logger.WithLogger(ctx, log)
		)

		if err := handler(ctxLog, msg, t); err != nil {
			log.WithError(err).Error("error while processing message")

			if _, err := botAPI.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: msg.ChatID,
				ReplyParameters: &models.ReplyParameters{
					MessageID: msg.MessageID,
				},
				Text: "internal error",
			}); err != nil {
				log.WithError(err).Error("send message error")
			}
		}
	}
}

func makeMsgFromUpdate(update *models.Update) BotMessage {
	if update.Message != nil {
		msg := BotMessage{
			MessageID: update.Message.ID,
			ChatID:    update.Message.Chat.ID,
			Text:      update.Message.Text,
		}

		if update.Message.SuccessfulPayment != nil {
			msg.Payment = Payment{
				IsSuccessful:   true,
				InvoicePayload: update.Message.SuccessfulPayment.InvoicePayload,
				Currency:       update.Message.SuccessfulPayment.Currency,
				TotalAmount:    update.Message.SuccessfulPayment.TotalAmount,
			}
		}

		return msg
	}

	if update.CallbackQuery != nil {
		if update.CallbackQuery.Message.Message != nil {
			return BotMessage{
				MessageID: update.CallbackQuery.Message.Message.ID,
				ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
				Data:      update.CallbackQuery.Data,
			}
		}

		if update.CallbackQuery.Message.InaccessibleMessage != nil {
			return BotMessage{
				MessageID: update.CallbackQuery.Message.InaccessibleMessage.MessageID,
				ChatID:    update.CallbackQuery.Message.InaccessibleMessage.Chat.ID,
				Data:      update.CallbackQuery.Data,
			}
		}
	}

	if update.PreCheckoutQuery != nil {
		return BotMessage{
			Payment: Payment{
				IsCheckout:     true,
				ID:             update.PreCheckoutQuery.ID,
				InvoicePayload: update.PreCheckoutQuery.InvoicePayload,
				Currency:       update.PreCheckoutQuery.Currency,
				TotalAmount:    update.PreCheckoutQuery.TotalAmount,
			},
		}
	}

	return BotMessage{}
}
