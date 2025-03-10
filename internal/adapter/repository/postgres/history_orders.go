package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (p *Postgres) HistoryOrders(ctx context.Context, chatID msginfo.ChatID, size int) ([]order.ShortOrder, error) {
	query, args, err := sqlx.Named(`
		SELECT
			id,
			chat_id,
			status,
			verification_code,
			currency_id,
			total_price,
			created_at,
			updated_at
		FROM
			orders
		WHERE
			chat_id = :chat_id
		ORDER BY
			id
		LIMIT
			:size
	`, map[string]any{
		"chat_id": chatID,
		"size":    size,
	})

	if err != nil {
		return nil, fmt.Errorf("sqlx named: %w", err)
	}

	var orders []model.Order
	if err := sqlx.SelectContext(ctx, p.db, &orders, p.db.Rebind(query), args...); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	portShortOrders, err := model.ToPortShortOrders(orders)
	if err != nil {
		return nil, fmt.Errorf("convert to port orders: %w", err)
	}

	return portShortOrders, nil
}
