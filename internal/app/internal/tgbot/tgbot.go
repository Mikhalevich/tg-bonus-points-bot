package tgbot

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/logger"
)

type TGBot struct {
	bot         *bot.Bot
	logger      logger.Logger
	middlewares []Middleware
	commands    []models.BotCommand
}

func New(b *bot.Bot, logger logger.Logger) *TGBot {
	return &TGBot{
		bot:    b,
		logger: logger,
	}
}
