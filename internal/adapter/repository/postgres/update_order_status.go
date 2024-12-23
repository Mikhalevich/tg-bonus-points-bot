package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/jmoiron/sqlx"
)

func (p *Postgres) UpdateOrderStatus(
	ctx context.Context,
	id order.ID,
	operationTime time.Time,
	newStatus order.Status,
	prevStatuses ...order.Status,
) error {
	query := fmt.Sprintf(`
		UPDATE orders SET
			status = :status,
			%s = :operation_time
		WHERE
			id = :id AND
			status IN (?)
	`, operationTimeFieldByStatus(newStatus))

	query, args, err := sqlx.Named(
		query,
		map[string]any{
			"status":         newStatus,
			"operation_time": operationTime,
			"id":             id.Int(),
		})

	if err != nil {
		return fmt.Errorf("named: %w", err)
	}

	args = append(args, prevStatuses)

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return fmt.Errorf("in statement %w", err)
	}

	res, err := p.db.ExecContext(ctx, p.db.Rebind(query), args...)
	if err != nil {
		return fmt.Errorf("exec context: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return errNotUpdated
	}

	return nil
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
