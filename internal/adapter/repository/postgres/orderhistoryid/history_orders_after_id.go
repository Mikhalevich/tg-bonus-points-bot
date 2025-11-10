package orderhistoryid

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (o *OrderHistoryID) HistoryOrdersAfterID(
	ctx context.Context,
	chatID msginfo.ChatID,
	afterOrderID order.ID,
	size int,
) ([]order.HistoryOrder, error) {
	query, args, err := sqlx.Named(`
		WITH orders_history AS (
			SELECT
				id,
				ROW_NUMBER() OVER (ORDER BY id) AS serial_number,
				status,
				currency_id,
				total_price,
				created_at
			FROM
				orders
			WHERE
				chat_id = :chat_id
		)
		SELECT
			id,
			serial_number,
			status,
			currency_id,
			total_price,
			created_at
		FROM
			orders_history
		WHERE
			id > :id
		ORDER BY
			id
		LIMIT
			:size
	`, map[string]any{
		"chat_id": chatID.Int64(),
		"id":      afterOrderID.Int(),
		"size":    size,
	})

	if err != nil {
		return nil, fmt.Errorf("sqlx named: %w", err)
	}

	orders, err := o.historyQuery(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("history query: %w", err)
	}

	return orders, nil
}
