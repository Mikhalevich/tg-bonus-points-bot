package setup

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	"github.com/Mikhalevich/tg-bonus-points-bot/cmd/manager/internal/app"
	"github.com/Mikhalevich/tg-bonus-points-bot/cmd/manager/internal/config"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/messagesender"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/driver"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/timeprovider"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/manager"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/logger"
)

func StartService(
	ctx context.Context,
	cfg config.Config,
	logger logger.Logger,
) error {
	botAPI, err := bot.New(cfg.Bot.Token, bot.WithSkipGetMe())
	if err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	pgDB, cleanup, err := MakePostgres(cfg.Postgres)
	if err != nil {
		return fmt.Errorf("make postgres: %w", err)
	}
	defer cleanup()

	var (
		sender           = messagesender.New(botAPI, cfg.Bot.PaymentToken)
		managerProcessor = manager.New(sender, pgDB, timeprovider.New())
	)

	if err := app.New(managerProcessor, logger).Start(
		ctx,
		cfg.HTTPPort,
	); err != nil {
		return fmt.Errorf("start bot: %w", err)
	}

	return nil
}

func MakePostgres(cfg config.Postgres) (*postgres.Postgres, func(), error) {
	if cfg.Connection == "" {
		return nil, func() {}, nil
	}

	driver := driver.NewPgx()

	dbConn, err := otelsql.Open(driver.Name(), cfg.Connection)
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}

	if err := dbConn.Ping(); err != nil {
		return nil, nil, fmt.Errorf("ping: %w", err)
	}

	p := postgres.New(dbConn, driver)

	return p, func() {
		dbConn.Close()
	}, nil
}
