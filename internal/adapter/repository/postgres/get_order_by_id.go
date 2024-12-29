package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (p *Postgres) GetOrderByID(ctx context.Context, id order.ID) (*order.Order, error) {
	var modelOrder model.Order
	if err := sqlx.GetContext(ctx, p.db, &modelOrder,
		`SELECT
			id,
			chat_id,
			status,
			verification_code,
			created_at,
			in_progress_at,
			ready_at,
			completed_at,
			canceled_at
		FROM
			orders
		WHERE
			id = $1
	`, id.Int(),
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}

		return nil, fmt.Errorf("get context: %w", err)
	}

	portOrder, err := model.ToPortOrder(&modelOrder)
	if err != nil {
		return nil, fmt.Errorf("convert to port order: %w", err)
	}

	return portOrder, nil
}
