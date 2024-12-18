package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mikhalevich/tg-loyalty-bot/internal/config"
	"github.com/Mikhalevich/tg-loyalty-bot/internal/infra"
	"github.com/Mikhalevich/tg-loyalty-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-loyalty-bot/internal/infra/tracing"
)

func main() {
	var cfg config.ScheduleBot
	if err := infra.LoadConfig(&cfg); err != nil {
		logger.StdLogger().WithError(err).Error("failed to load config")
		os.Exit(1)
	}

	log, err := infra.SetupLogger(cfg.LogLevel)
	if err != nil {
		logger.StdLogger().WithError(err).Error("failed to setup logger")
		os.Exit(1)
	}

	if err := runService(cfg, log); err != nil {
		log.WithError(err).Error("failed run service")
		os.Exit(1)
	}
}

func runService(cfg config.ScheduleBot, log logger.Logger) error {
	if err := tracing.SetupTracer(cfg.Tracing.Endpoint, cfg.Tracing.ServiceName, ""); err != nil {
		return fmt.Errorf("setup tracer: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Info("starting bot")

	if err := infra.StartBot(
		ctx,
		cfg.Bot.Token,
		cfg.Postgres,
		log.WithField("bot_name", "schedule"),
	); err != nil {
		return fmt.Errorf("start bot: %w", err)
	}

	log.Info("bot stopped")

	return nil
}
