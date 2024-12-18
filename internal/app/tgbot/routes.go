package tgbot

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/tgbot/router"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

type OrderProcessor interface {
	MakeOrder(ctx context.Context, msgInfo msginfo.Info) error
}

func Routes(o OrderProcessor) RouteRegisterFunc {
	return func(r router.Register) {
		r.AddExactTextRoute("/start", makeStartHandler(r))

		r.AddExactTextRoute("/order", o.MakeOrder)

		r.AddDefaultTextHandler(makeDefaultHandler(r))
	}
}

func makeStartHandler(r router.Register) msginfo.Handler {
	return func(ctx context.Context, info msginfo.Info) error {
		if err := r.SendMessage(
			ctx,
			info.ChatID.Int64(),
			"Type /makeorder for requesting order",
		); err != nil {
			return fmt.Errorf("send message: %w", err)
		}

		return nil
	}
}

func makeDefaultHandler(r router.Register) msginfo.Handler {
	return func(ctx context.Context, info msginfo.Info) error {
		// skip message.
		return nil
	}
}
