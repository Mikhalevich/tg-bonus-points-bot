package tgbot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

type Probe func(ctx context.Context) error

type TGBot struct {
	bot              *bot.Bot
	isWebHook        bool
	logger           logger.Logger
	middlewares      []Middleware
	commands         []models.BotCommand
	defaultHandlerFn Handler
	livenessProbe    Probe
	readinessProbe   Probe
}

func New(
	token string,
	webHookToken string,
	logger logger.Logger,
) (*TGBot, error) {
	tgbot := &TGBot{
		isWebHook: webHookToken != "",
		logger:    logger,
	}

	opts := []bot.Option{
		bot.WithSkipGetMe(),
		bot.WithDefaultHandler(tgbot.makeDefaultHandler()),
	}

	if webHookToken != "" {
		opts = append(opts, bot.WithWebhookSecretToken(webHookToken))
	}

	botAPI, err := bot.New(
		token,
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("creating bot: %w", err)
	}

	tgbot.bot = botAPI

	return tgbot, nil
}
