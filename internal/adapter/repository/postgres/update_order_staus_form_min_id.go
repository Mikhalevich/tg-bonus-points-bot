package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/transaction"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (p *Postgres) UpdateOrderStatusForMinID(
	ctx context.Context,
	operationTime time.Time,
	newStatus, prevStatus order.Status,
) (*order.Order, error) {
	var (
		dbOrder       *model.Order
		orderProducts []model.OrderProductFull
		timeline      []model.OrderTimeline
		err           error
	)

	if err := transaction.Transaction(ctx, p.db, true,
		func(ctx context.Context, tx sqlx.ExtContext) error {
			dbOrder, err = updateOrderStatusForMinID(ctx, tx, newStatus, prevStatus)
			if err != nil {
				return fmt.Errorf("update order status for min id: %w", err)
			}

			if err := insertOrderTimeline(ctx, tx, model.OrderTimeline{
				ID:        dbOrder.ID,
				Status:    newStatus.String(),
				UpdatedAt: operationTime,
			}); err != nil {
				return fmt.Errorf("insert order timeline: %w", err)
			}

			orderProducts, err = selectOrderProducts(ctx, p.db, dbOrder.ID)
			if err != nil {
				return fmt.Errorf("select order products: %w", err)
			}

			timeline, err = selectOrderTimeline(ctx, tx, dbOrder.ID)
			if err != nil {
				return fmt.Errorf("select order timeline: %w", err)
			}

			return nil
		},
	); err != nil {
		return nil, fmt.Errorf("transaction: %w", err)
	}

	portOrder, err := model.ToPortOrder(dbOrder, orderProducts, timeline)
	if err != nil {
		return nil, fmt.Errorf("convert to port order: %w", err)
	}

	return portOrder, nil
}

func updateOrderStatusForMinID(
	ctx context.Context,
	ext sqlx.ExtContext,
	newStatus order.Status,
	prevStatus order.Status,
) (*model.Order, error) {
	query, args, err := sqlx.Named(`
		UPDATE orders SET
			status = :new_status
		WHERE id = (
				SELECT MIN(id)
				FROM orders
				WHERE status = :previous_status
			)
		RETURNING *
		`, map[string]any{
		"new_status":      newStatus,
		"previous_status": prevStatus,
	})

	if err != nil {
		return nil, fmt.Errorf("named: %w", err)
	}

	var dbOrder model.Order
	if err := sqlx.GetContext(ctx, ext, &dbOrder, ext.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}

		return nil, fmt.Errorf("get context: %w", err)
	}

	return &dbOrder, nil
}
