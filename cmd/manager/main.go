package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/config"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/tracing"
)

func main() {
	var cfg config.ManagerHTTPService
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

func runService(cfg config.ManagerHTTPService, log logger.Logger) error {
	if err := tracing.SetupTracer(cfg.Tracing.Endpoint, cfg.Tracing.ServiceName, ""); err != nil {
		return fmt.Errorf("setup tracer: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Info("service starting")

	if err := infra.StartManagerService(
		ctx,
		cfg.HTTPPort,
		cfg.Bot,
		cfg.Postgres,
		log.WithField("service_name", "http_manager"),
	); err != nil {
		return fmt.Errorf("start bot: %w", err)
	}

	log.Info("service stopped")

	return nil
}
