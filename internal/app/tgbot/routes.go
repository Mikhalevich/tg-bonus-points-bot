package tgbot

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/tgbot/router"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

type OrderProcessor interface {
	MakeOrder(ctx context.Context, msgInfo msginfo.Info) error
	GetActiveOrder(ctx context.Context, msgInfo msginfo.Info) error
}

func Routes(o OrderProcessor) RouteRegisterFunc {
	return func(r router.Register) {
		r.AddTextCommand("/start", makeStartHandler(r))

		r.AddMenuCommand("/order", "make order", o.MakeOrder)
		r.AddMenuCommand("/get_active_order", "retrieve active order", o.GetActiveOrder)

		r.AddDefaultTextHandler(makeDefaultHandler(r))
	}
}

func makeStartHandler(r router.Register) msginfo.Handler {
	return func(ctx context.Context, info msginfo.Info) error {
		if err := r.SendMessage(
			ctx,
			info.ChatID.Int64(),
			"Type /order for requesting an order",
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
