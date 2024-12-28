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

func (p *Postgres) UpdateOrderStatus(
	ctx context.Context,
	id order.ID,
	operationTime time.Time,
	newStatus order.Status,
	prevStatuses ...order.Status,
) (*order.Order, error) {
	query := fmt.Sprintf(`
		UPDATE orders SET
			status = :status,
			%s = :operation_time
		WHERE
			id = :id AND
			status IN (?)
		RETURNING *
	`, operationTimeFieldByStatus(newStatus))

	query, args, err := sqlx.Named(
		query,
		map[string]any{
			"status":         newStatus,
			"operation_time": operationTime,
			"id":             id.Int(),
		})

	if err != nil {
		return nil, fmt.Errorf("named: %w", err)
	}

	args = append(args, prevStatuses)

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, fmt.Errorf("in statement %w", err)
	}

	var modOrder model.Order
	if err := sqlx.GetContext(ctx, p.db, &modOrder, p.db.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotUpdated
		}

		return nil, fmt.Errorf("exec context: %w", err)
	}

	portOrder, err := model.ToPortOrder(&modOrder)
	if err != nil {
		return nil, fmt.Errorf("convert to port order: %w", err)
	}

	return portOrder, nil
}

func operationTimeFieldByStatus(s order.Status) string {
	switch s {
	case order.StatusCreated:
		return "created_at"
	case order.StatusInProgress:
		return "in_progress_at"
	case order.StatusReady:
		return "ready_at"
	case order.StatusCompleted:
		return "completed_at"
	case order.StatusCanceled:
		return "canceled_at"
	}

	return "invalid"
}
