package infra

import (
	"context"
	"flag"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/jinzhu/configor"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/messagesender"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/qrcodegenerator"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/driver"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/httpmanager"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/config"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/customer"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/manager"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/logger"
)

func LoadConfig(cfg any) error {
	configFile := flag.String("config", "config/config.yaml", "consumer worker config file")
	flag.Parse()

	if err := configor.Load(cfg, *configFile); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return nil
}

func SetupLogger(lvl string) (logger.Logger, error) {
	log, err := logger.NewLogrusWithLevel(lvl)
	if err != nil {
		return nil, fmt.Errorf("creating new logger: %w", err)
	}

	logger.SetStdLogger(log)

	return log, nil
}

func StartBot(
	ctx context.Context,
	botAPItoken string,
	postgresCfg config.Postgres,
	logger logger.Logger,
) error {
	b, err := bot.New(botAPItoken, bot.WithSkipGetMe())
	if err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	pg, cleanup, err := MakePostgres(postgresCfg)
	if err != nil {
		return fmt.Errorf("make postgres: %w", err)
	}
	defer cleanup()

	var (
		sender            = messagesender.New(b)
		qrGenerator       = qrcodegenerator.New()
		customerProcessor = customer.New(sender, qrGenerator, pg)
	)

	if err := tgbot.Start(
		ctx,
		b,
		logger,
		tgbot.Routes(customerProcessor),
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

	db, err := otelsql.Open(driver.Name(), cfg.Connection)
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("ping: %w", err)
	}

	p := postgres.New(db, driver)

	return p, func() {
		db.Close()
	}, nil
}

func StartManagerService(
	ctx context.Context,
	httpPort int,
	botAPItoken string,
	postgresCfg config.Postgres,
	logger logger.Logger,
) error {
	b, err := bot.New(botAPItoken, bot.WithSkipGetMe())
	if err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	pg, cleanup, err := MakePostgres(postgresCfg)
	if err != nil {
		return fmt.Errorf("make postgres: %w", err)
	}
	defer cleanup()

	var (
		sender           = messagesender.New(b)
		managerProcessor = manager.New(sender, pg)
	)

	if err := httpmanager.New(managerProcessor, logger).Start(
		ctx,
		httpPort,
	); err != nil {
		return fmt.Errorf("start bot: %w", err)
	}

	return nil
}
