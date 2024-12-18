package tgbot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/tgbot/router"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/logger"
)

type RouteRegisterFunc func(register router.Register)

func Start(
	ctx context.Context,
	b *bot.Bot,
	logger logger.Logger,
	routesRegisterFn RouteRegisterFunc,
) error {
	r := router.New(b, logger)

	routesRegisterFn(r)

	if err := r.SetMyCommands(ctx); err != nil {
		return fmt.Errorf("set my commands: %w", err)
	}

	b.Start(ctx)

	return nil
}
