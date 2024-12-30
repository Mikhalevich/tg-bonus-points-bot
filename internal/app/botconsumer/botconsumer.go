package botconsumer

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/botconsumer/tghandler"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/internal/tgbot"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/logger"
)

func Start(
	ctx context.Context,
	b *bot.Bot,
	logger logger.Logger,
	processor tghandler.OrderProcessor,
) error {
	var (
		tbot       = tgbot.New(b, logger)
		botHandler = tghandler.New(processor)
	)

	makeRoutes(tbot, botHandler)

	if err := tbot.SetMyCommands(ctx); err != nil {
		return fmt.Errorf("set my commands: %w", err)
	}

	b.Start(ctx)

	return nil
}
