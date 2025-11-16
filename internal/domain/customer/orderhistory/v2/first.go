package v2

import (
	"context"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

func (oh *OrderHistory) First(ctx context.Context, info msginfo.Info) error {
	return oh.Page(ctx, info, 1)
}
