package botconsumer

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/botconsumer/tghandler"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/internal/tgbot"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/logger"
)

func Start(
	ctx context.Context,
	token string,
	logger logger.Logger,
	cartProcessor tghandler.CartProcessor,
	actionProcessor tghandler.OrderActionProcessor,
	paymentProcessor tghandler.OrderPaymentProcessor,
	buttonProvider tghandler.ButtonProvider,
) error {
	var (
		botHandler = tghandler.New(cartProcessor, actionProcessor, paymentProcessor, buttonProvider)
	)

	tbot, err := tgbot.New(token, logger)
	if err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	makeRoutes(tbot, botHandler)

	if err := tbot.Start(ctx); err != nil {
		return fmt.Errorf("bot start: %w", err)
	}

	return nil
}
