package infra

import (
	"context"
	"flag"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/jinzhu/configor"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/buttonrespository"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/cart"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/messagesender"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/qrcodegenerator"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/driver"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/botconsumer"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/httpmanager"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/config"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/customer"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/manager"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
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
	botCfg config.Bot,
	postgresCfg config.Postgres,
	cartRedisCfg config.CartRedis,
	buttonRedisCfg config.ButtonRedis,
	logger logger.Logger,
) error {
	b, err := bot.New(botCfg.Token, bot.WithSkipGetMe())
	if err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	pg, cleanup, err := MakePostgres(postgresCfg)
	if err != nil {
		return fmt.Errorf("make postgres: %w", err)
	}
	defer cleanup()

	cartRedis, err := MakeRedisCart(ctx, cartRedisCfg)
	if err != nil {
		return fmt.Errorf("make redis cart: %w", err)
	}

	buttonRepository, err := MakeRedisButtonRepository(ctx, buttonRedisCfg)
	if err != nil {
		return fmt.Errorf("make redis button repository: %w", err)
	}

	var (
		sender            = messagesender.New(b, botCfg.PaymentToken)
		qrGenerator       = qrcodegenerator.New()
		customerProcessor = customer.New(sender, qrGenerator, pg, cartRedis, buttonRepository)
	)

	if err := botconsumer.Start(
		ctx,
		botCfg.Token,
		logger,
		customerProcessor,
	); err != nil {
		return fmt.Errorf("start bot: %w", err)
	}

	return nil
}

func MakeRedisButtonRepository(ctx context.Context, cfg config.ButtonRedis) (port.ButtonRepository, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Pwd,
		DB:       cfg.DB,
	})

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		return nil, fmt.Errorf("redis instrument tracing: %w", err)
	}

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return buttonrespository.New(rdb, cfg.TTL), nil
}

func MakeRedisCart(ctx context.Context, cfg config.CartRedis) (port.Cart, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Pwd,
		DB:       cfg.DB,
	})

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		return nil, fmt.Errorf("redis instrument tracing: %w", err)
	}

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return cart.New(rdb, cfg.TTL), nil
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
	botCfg config.Bot,
	postgresCfg config.Postgres,
	logger logger.Logger,
) error {
	b, err := bot.New(botCfg.Token, bot.WithSkipGetMe())
	if err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	pg, cleanup, err := MakePostgres(postgresCfg)
	if err != nil {
		return fmt.Errorf("make postgres: %w", err)
	}
	defer cleanup()

	var (
		sender           = messagesender.New(b, botCfg.PaymentToken)
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
