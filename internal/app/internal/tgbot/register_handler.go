package tgbot

import (
	"context"
	"strings"

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
		bot.MatchTypePrefix,
		t.wrapHandler(command, command, handler),
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
		t.wrapHandler("default_text_handler", "", h),
	)
}

func (t *TGBot) AddDefaultCallbackQueryHander(h Handler) {
	t.bot.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"",
		bot.MatchTypePrefix,
		t.wrapHandler("default_callback_query", "", h),
	)
}

func (t *TGBot) wrapHandler(pattern, prefixTrimFromText string, h Handler) bot.HandlerFunc {
	h = t.applyMiddleware(h)

	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		ctx, span := tracing.StartSpanName(ctx, pattern)
		defer span.End()

		var (
			msg    = makeMsgFromUpdate(update, prefixTrimFromText)
			log    = t.logger.WithContext(ctx).WithField("endpoint", pattern)
			ctxLog = logger.WithLogger(ctx, log)
		)

		if err := h(ctxLog, msg, t); err != nil {
			log.WithError(err).Error("error while processing message")

			if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
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

func makeMsgFromUpdate(u *models.Update, prefixTrimFromText string) BotMessage {
	if u.Message != nil {
		msg := BotMessage{
			MessageID: u.Message.ID,
			ChatID:    u.Message.Chat.ID,
			Text:      strings.TrimSpace(strings.TrimPrefix(u.Message.Text, prefixTrimFromText)),
		}

		if u.Message.SuccessfulPayment != nil {
			msg.Payment = Payment{
				IsSuccessful:   true,
				InvoicePayload: u.Message.SuccessfulPayment.InvoicePayload,
				Currency:       u.Message.SuccessfulPayment.Currency,
				TotalAmount:    u.Message.SuccessfulPayment.TotalAmount,
			}
		}

		return msg
	}

	if u.CallbackQuery != nil {
		if u.CallbackQuery.Message.Message != nil {
			return BotMessage{
				MessageID: u.CallbackQuery.Message.Message.ID,
				ChatID:    u.CallbackQuery.Message.Message.Chat.ID,
				Data:      u.CallbackQuery.Data,
			}
		}

		if u.CallbackQuery.Message.InaccessibleMessage != nil {
			return BotMessage{
				MessageID: u.CallbackQuery.Message.InaccessibleMessage.MessageID,
				ChatID:    u.CallbackQuery.Message.InaccessibleMessage.Chat.ID,
				Data:      u.CallbackQuery.Data,
			}
		}
	}

	if u.PreCheckoutQuery != nil {
		return BotMessage{
			Payment: Payment{
				IsCheckout:     true,
				ID:             u.PreCheckoutQuery.ID,
				InvoicePayload: u.PreCheckoutQuery.InvoicePayload,
				Currency:       u.PreCheckoutQuery.Currency,
				TotalAmount:    u.PreCheckoutQuery.TotalAmount,
			},
		}
	}

	return BotMessage{}
}
