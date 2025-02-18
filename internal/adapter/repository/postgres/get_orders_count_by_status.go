package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (p *Postgres) GetOrdersCountByStatus(ctx context.Context, statuses ...order.Status) (int, error) {
	query, args, err := sqlx.In(`
		SELECT
			COUNT(*)
		FROM
			orders
		WHERE
			status IN(?)
	`, statuses)

	if err != nil {
		return 0, fmt.Errorf("sqlx in statement: %w", err)
	}

	var count int
	if err := sqlx.GetContext(ctx, p.db, &count, p.db.Rebind(query), args...); err != nil {
		return 0, fmt.Errorf("get context: %w", err)
	}

	return count, nil
}
