package setup

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/msgconsumer/internal/app"
	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/msgconsumer/internal/app/kafkaconsumer"
	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/msgconsumer/internal/config"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

func StartConsumer(
	ctx context.Context,
	cfg config.Config,
	logger logger.Logger,
) error {
	var (
		consumer = kafkaconsumer.New(cfg.Kafka)
	)

	if err := app.StartApp(ctx, consumer); err != nil {
		return fmt.Errorf("start app: %w", err)
	}

	return nil
}
