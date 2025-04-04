package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (p *Postgres) GetOrderPositionByStatus(
	ctx context.Context,
	orderID order.ID,
	statuses ...order.Status,
) (int, error) {
	query, args, err := sqlx.In(`
		WITH order_queue AS (
			SELECT
				id,
				ROW_NUMBER() OVER (ORDER BY updated_at::date, daily_position) AS position
			FROM
				orders
			WHERE
				status IN(?)
			ORDER BY
				updated_at::date,
				daily_position
		)
		SELECT
			position
		FROM
			order_queue
		WHERE
			id = ?
	`, statuses, orderID.Int())

	if err != nil {
		return 0, fmt.Errorf("sqlx in statement: %w", err)
	}

	var pos int
	if err := sqlx.GetContext(ctx, p.db, &pos, p.db.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errNotFound
		}

		return 0, fmt.Errorf("get context: %w", err)
	}

	return pos, nil
}
