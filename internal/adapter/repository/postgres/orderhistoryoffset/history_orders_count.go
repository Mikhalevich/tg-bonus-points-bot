package orderhistoryoffset

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (oh *OrderHistoryOffset) HistoryOrdersCount(
	ctx context.Context,
	chatID msginfo.ChatID,
) (int, error) {
	query, args, err := sqlx.Named(`
		SELECT
			COUNT(*)
		FROM
			orders
		WHERE
			chat_id = :chat_id
	`,
		map[string]any{
			"chat_id": chatID.Int64(),
		},
	)

	if err != nil {
		return 0, fmt.Errorf("prepare query: %w", err)
	}

	var count int
	if err := sqlx.GetContext(ctx, oh.db, &count, oh.db.Rebind(query), args...); err != nil {
		return 0, fmt.Errorf("get context: %w", err)
	}

	return count, nil
}
