package v2

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (oh *OrderHistory) Show(ctx context.Context, info msginfo.Info) error {
	if err := oh.loadPageByNumber(ctx, info, 1, SendMessage); err != nil {
		return fmt.Errorf("load page by number: %w", err)
	}

	return nil
}
