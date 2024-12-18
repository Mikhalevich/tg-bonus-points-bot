package router

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-loyalty-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-loyalty-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-loyalty-bot/internal/infra/tracing"
)

type Middleware func(next msginfo.Handler) msginfo.Handler

type Register interface {
	AddExactTextRoute(pattern string, handler msginfo.Handler)
	AddDefaultTextHandler(handler msginfo.Handler)
	AddDefaultCallbackQueryHander(h msginfo.Handler)
	AddMiddleware(middleware Middleware)
	MiddlewareGroup(fn func(r Register))
	SendMessage(ctx context.Context, chatID int64, msg string) error
}

type Router struct {
	bot         *bot.Bot
	logger      logger.Logger
	middlewares []Middleware
}

func New(b *bot.Bot, logger logger.Logger) *Router {
	return &Router{
		bot:    b,
		logger: logger,
	}
}

func (r *Router) SendMessage(ctx context.Context, chatID int64, msg string) error {
	if _, err := r.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   msg,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (r *Router) MiddlewareGroup(fn func(r Register)) {
	group := &Router{
		bot:         r.bot,
		logger:      r.logger,
		middlewares: r.middlewares[:len(r.middlewares):len(r.middlewares)],
	}

	fn(group)
}

func (r *Router) AddExactTextRoute(pattern string, handler msginfo.Handler) {
	r.bot.RegisterHandler(
		bot.HandlerTypeMessageText,
		pattern,
		bot.MatchTypeExact,
		r.wrapHandler(pattern, handler),
	)
}

func (r *Router) AddDefaultTextHandler(h msginfo.Handler) {
	r.bot.RegisterHandler(
		bot.HandlerTypeMessageText,
		"",
		bot.MatchTypePrefix,
		r.wrapHandler("default_text_handler", h),
	)
}

func (r *Router) AddDefaultCallbackQueryHander(h msginfo.Handler) {
	r.bot.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"",
		bot.MatchTypePrefix,
		r.wrapHandler("default_callback_query", h),
	)
}

func (r *Router) AddMiddleware(m Middleware) {
	r.middlewares = append(r.middlewares, m)
}

func (r *Router) wrapHandler(pattern string, h msginfo.Handler) bot.HandlerFunc {
	h = r.applyMiddleware(h)

	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		ctx, span := tracing.StartSpanName(ctx, pattern)
		defer span.End()

		var (
			msgInfo = makeMsgInfoFromUpdate(update)
			log     = r.logger.WithContext(ctx).WithField("endpoint", pattern)
			ctxLog  = logger.WithLogger(ctx, log)
		)

		if err := h(ctxLog, msgInfo); err != nil {
			log.WithError(err).Error("error while processing message")

			if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: msgInfo.ChatID,
				ReplyParameters: &models.ReplyParameters{
					MessageID: msgInfo.MessageID.Int(),
				},
				Text: "internal error",
			}); err != nil {
				log.WithError(err).Error("send internal error message")
			}
		}
	}
}

func (r *Router) applyMiddleware(h msginfo.Handler) msginfo.Handler {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		h = r.middlewares[i](h)
	}

	return h
}

func makeMsgInfoFromUpdate(u *models.Update) msginfo.Info {
	if u.Message != nil {
		return msginfo.Info{
			MessageID: msginfo.MessageIDFromInt(u.Message.ID),
			ChatID:    msginfo.ChatIDFromInt(u.Message.Chat.ID),
			Text:      u.Message.Text,
		}
	}

	if u.CallbackQuery != nil {
		if u.CallbackQuery.Message.Message != nil {
			return msginfo.Info{
				MessageID: msginfo.MessageIDFromInt(u.CallbackQuery.Message.Message.ID),
				ChatID:    msginfo.ChatIDFromInt(u.CallbackQuery.Message.Message.Chat.ID),
				Data:      u.CallbackQuery.Data,
			}
		}

		if u.CallbackQuery.Message.InaccessibleMessage != nil {
			return msginfo.Info{
				MessageID: msginfo.MessageIDFromInt(u.CallbackQuery.Message.InaccessibleMessage.MessageID),
				ChatID:    msginfo.ChatIDFromInt(u.CallbackQuery.Message.InaccessibleMessage.Chat.ID),
				Data:      u.CallbackQuery.Data,
			}
		}
	}

	return msginfo.Info{}
}
