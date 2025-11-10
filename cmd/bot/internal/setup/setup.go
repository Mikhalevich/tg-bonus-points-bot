package setup

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	"github.com/Mikhalevich/tg-bonus-points-bot/cmd/bot/internal/app"
	"github.com/Mikhalevich/tg-bonus-points-bot/cmd/bot/internal/config"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/buttonrespository"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/cart"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/dailypositiongenerator"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/messagesender"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/qrcodegenerator"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/driver"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/timeprovider"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/verificationcodegenerator"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/buttonprovider"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/customer/cartprocessing"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/customer/orderaction"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/customer/orderhistory"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/customer/orderpayment"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/logger"
)

func StartBot(
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

	cartRedis, err := MakeRedisCart(ctx, cfg.CartRedis)
	if err != nil {
		return fmt.Errorf("make redis cart: %w", err)
	}

	dailyPosition, err := MakeRedisDailyPositionGenerator(ctx, cfg.DailyPositionRedis)
	if err != nil {
		return fmt.Errorf("make redis daily position generator: %w", err)
	}

	buttonRepository, err := MakeRedisButtonRepository(ctx, cfg.ButtonRedis)
	if err != nil {
		return fmt.Errorf("make redis button repository: %w", err)
	}

	var (
		sender        = messagesender.New(botAPI, cfg.Bot.PaymentToken)
		qrGenerator   = qrcodegenerator.New()
		cartProcessor = cartprocessing.New(cfg.StoreID, pgDB, pgDB, cartRedis, sender,
			timeprovider.New(), buttonRepository)
		actionProcessor  = orderaction.New(sender, pgDB, buttonRepository, timeprovider.New())
		historyProcessor = orderhistory.New(pgDB, pgDB, sender, buttonRepository, cfg.OrderHistory.PageSize)
		paymentProcessor = orderpayment.New(cfg.StoreID, sender, qrGenerator, pgDB, pgDB,
			dailyPosition, verificationcodegenerator.New(), timeprovider.New())
		buttonProvider = buttonprovider.New(buttonRepository)
	)

	if err := app.Start(
		ctx,
		cfg.Bot.Token,
		logger,
		cartProcessor,
		actionProcessor,
		historyProcessor,
		paymentProcessor,
		buttonProvider,
	); err != nil {
		return fmt.Errorf("start bot: %w", err)
	}

	return nil
}

func MakeRedisButtonRepository(
	ctx context.Context,
	cfg config.ButtonRedis,
) (*buttonrespository.ButtonRepository, error) {
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

func MakeRedisDailyPositionGenerator(
	ctx context.Context,
	cfg config.DailyPositionRedis,
) (port.DailyPositionGenerator, error) {
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

	return dailypositiongenerator.New(rdb, cfg.TTL), nil
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
