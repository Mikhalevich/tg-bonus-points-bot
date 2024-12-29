package tgbot

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/tgbot/router"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

type OrderProcessor interface {
	MakeOrder(ctx context.Context, msgInfo msginfo.Info) error
	GetActiveOrder(ctx context.Context, msgInfo msginfo.Info) error
	CancelOrder(ctx context.Context, id order.ID) error
}

func Routes(o OrderProcessor) RouteRegisterFunc {
	return func(r router.Register) {
		r.AddTextCommand("/start", makeStartHandler(r))

		r.AddMenuCommand("/order", "make order", o.MakeOrder)
		r.AddMenuCommand("/get_active_order", "retrieve active order", o.GetActiveOrder)

		r.AddDefaultTextHandler(makeDefaultHandler(r))
		r.AddDefaultCallbackQueryHander(makeDefaultCallbackQueryHandler(o))
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

func makeDefaultCallbackQueryHandler(processor OrderProcessor) msginfo.Handler {
	return func(ctx context.Context, info msginfo.Info) error {
		orderID, err := order.IDFromString(info.Data)
		if err != nil {
			return fmt.Errorf("invalid order id: %w", err)
		}

		if err := processor.CancelOrder(ctx, orderID); err != nil {
			return fmt.Errorf("cancel order: %w", err)
		}

		return nil
	}
}
