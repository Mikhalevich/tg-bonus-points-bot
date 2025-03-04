package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/config"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/tracing"
)

func main() {
	var cfg config.ConsumerBot
	if err := infra.LoadConfig(&cfg); err != nil {
		logger.StdLogger().WithError(err).Error("failed to load config")
		os.Exit(1)
	}

	log, err := infra.SetupLogger(cfg.LogLevel)
	if err != nil {
		logger.StdLogger().WithError(err).Error("failed to setup logger")
		os.Exit(1)
	}

	if err := tracing.SetupTracer(cfg.Tracing.Endpoint, cfg.Tracing.ServiceName, ""); err != nil {
		log.WithError(err).Error("failed to setup tracer")
		os.Exit(1)
	}

	if err := infra.RunSignalInterruptionFunc(func(ctx context.Context) error {
		log.Info("starting consumer bot")

		if err := infra.StartBot(
			ctx,
			cfg.StoreID,
			cfg.Bot,
			cfg.Postgres,
			cfg.CartRedis,
			cfg.DailyPositionRedis,
			cfg.ButtonRedis,
			log.WithField("bot_name", "consumer"),
		); err != nil {
			return fmt.Errorf("start bot: %w", err)
		}

		log.Info("consumer bot stopped")

		return nil
	}); err != nil {
		log.WithError(err).Error("failed run service")
		os.Exit(1)
	}
}
