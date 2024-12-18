package order

import (
	"context"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

func (o *Order) MakeOrder(ctx context.Context, info msginfo.Info) error {
	o.sender.ReplyText(ctx, info.ChatID, info.MessageID, "make order stub")

	return nil
}
