package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (p *Postgres) HistoryOrders(ctx context.Context, chatID msginfo.ChatID, size int) ([]order.HistoryOrder, error) {
	query, args, err := sqlx.Named(`
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
		ORDER BY
			id DESC
		LIMIT
			:size
	`, map[string]any{
		"chat_id": chatID,
		"size":    size,
	})

	if err != nil {
		return nil, fmt.Errorf("sqlx named: %w", err)
	}

	var orders []model.HistoryOrder
	if err := sqlx.SelectContext(ctx, p.db, &orders, p.db.Rebind(query), args...); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	portHistoryOrders, err := model.ToPortHistoryOrders(orders)
	if err != nil {
		return nil, fmt.Errorf("convert to port orders: %w", err)
	}

	return portHistoryOrders, nil
}
