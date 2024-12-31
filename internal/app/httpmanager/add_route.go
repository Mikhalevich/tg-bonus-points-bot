package httpmanager

import (
	"context"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/tracing"
)

type handlerFunc[I, O any] func(context.Context, *I) (*O, error)

func addRoute[I, O any](m *HTTPManager, op huma.Operation, hf handlerFunc[I, O]) {
	huma.Register(
		m.humaAPI,
		op,
		makeHandlerWrapper(m, op.Path, hf),
	)
}

func makeHandlerWrapper[I, O any](m *HTTPManager, pattern string, hf handlerFunc[I, O]) handlerFunc[I, O] {
	return func(ctx context.Context, input *I) (*O, error) {
		ctx, span := tracing.StartSpanName(ctx, pattern)
		defer span.End()

		var (
			log    = m.logger.WithContext(ctx).WithField("endpoint", pattern)
			ctxLog = logger.WithLogger(ctx, log)
		)

		output, err := hf(ctxLog, input)
		if err != nil {
			log.WithError(err).
				Error("http handler error")
		}

		return output, err
	}
}
