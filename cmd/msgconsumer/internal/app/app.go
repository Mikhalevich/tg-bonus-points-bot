package app

import (
	"context"
	"fmt"
)

type Consumer interface {
	Consume(
		ctx context.Context,
		processFn func(ctx context.Context, payload []byte) error,
	) error
}

func StartApp(
	ctx context.Context,
	consumer Consumer,
) error {
	if err := consumer.Consume(ctx, func(ctx context.Context, payload []byte) error {
		return nil
	}); err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	return nil
}
