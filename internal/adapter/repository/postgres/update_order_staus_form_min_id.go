package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (p *Postgres) UpdateOrderStatusForMinID(
	ctx context.Context,
	operationTime time.Time,
	prevStatus, newStatus order.Status,
) (*order.Order, error) {
	query := fmt.Sprintf(`
		UPDATE orders SET
			status = :new_status,
			%s = :operation_time
		WHERE id = (
				SELECT MIN(id)
				FROM orders
				WHERE status = :previous_status
			)
		RETURNING *
	`, operationTimeFieldByStatus(newStatus))

	query, args, err := sqlx.Named(
		query,
		map[string]any{
			"new_status":      newStatus,
			"operation_time":  operationTime,
			"previous_status": prevStatus,
		})

	if err != nil {
		return nil, fmt.Errorf("named: %w", err)
	}

	var modelOrder model.Order
	if err := sqlx.GetContext(ctx, p.db, &modelOrder, p.db.Rebind(query), args...); err != nil {
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
