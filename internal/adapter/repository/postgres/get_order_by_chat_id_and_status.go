package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (p *Postgres) GetOrderByChatIDAndStatus(
	ctx context.Context,
	id msginfo.ChatID,
	statuses ...order.Status,
) (*order.Order, error) {
	query, args, err := sqlx.In(`
		SELECT
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
			chat_id = ? AND
			status IN(?)
	`, id, statuses)

	if err != nil {
		return nil, fmt.Errorf("sqlx.in: %w", err)
	}

	var order model.Order
	if err := sqlx.GetContext(ctx, p.db, &order, p.db.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}

		return nil, fmt.Errorf("get context: %w", err)
	}

	portOrder, err := model.ToPortOrder(&order)
	if err != nil {
		return nil, fmt.Errorf("convert to port order: %w", err)
	}

	return portOrder, nil
}
