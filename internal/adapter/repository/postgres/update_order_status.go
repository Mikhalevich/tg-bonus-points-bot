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

func (p *Postgres) UpdateOrderStatus(
	ctx context.Context,
	id order.ID,
	operationTime time.Time,
	newStatus order.Status,
	prevStatuses ...order.Status,
) (*order.Order, error) {
	var (
		dbOrder       *model.Order
		orderProducts []model.OrderProductFull
		orderTimeline []model.OrderTimeline
		err           error
	)

	if err := transaction.Transaction(ctx, p.db, true,
		func(ctx context.Context, tx sqlx.ExtContext) error {
			dbOrder, err = updateOrderStatus(ctx, tx, id, newStatus, prevStatuses...)
			if err != nil {
				return fmt.Errorf("update order status: %w", err)
			}

			if err := insertOrderTimeline(ctx, tx, model.OrderTimeline{
				ID:        id.Int(),
				Status:    newStatus.String(),
				UpdatedAt: operationTime,
			}); err != nil {
				return fmt.Errorf("insert order timeline: %w", err)
			}

			orderProducts, err = selectOrderProducts(ctx, p.db, dbOrder.ID)
			if err != nil {
				return fmt.Errorf("select order products: %w", err)
			}

			orderTimeline, err = selectOrderTimeline(ctx, tx, id.Int())
			if err != nil {
				return fmt.Errorf("select order timeline: %w", err)
			}

			return nil
		},
	); err != nil {
		return nil, fmt.Errorf("transaction: %w", err)
	}

	portOrder, err := model.ToPortOrder(dbOrder, orderProducts, orderTimeline)
	if err != nil {
		return nil, fmt.Errorf("convert to port order: %w", err)
	}

	return portOrder, nil
}

func updateOrderStatus(
	ctx context.Context,
	ext sqlx.ExtContext,
	id order.ID,
	newStatus order.Status,
	prevStatuses ...order.Status,
) (*model.Order, error) {
	query, args, err := sqlx.Named(`
		UPDATE orders SET
			status = :status
		WHERE
			id = :id AND
			status IN (?)
		RETURNING *
		`,
		map[string]any{
			"status": newStatus,
			"id":     id.Int(),
		})

	if err != nil {
		return nil, fmt.Errorf("named: %w", err)
	}

	args = append(args, prevStatuses)

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, fmt.Errorf("in statement %w", err)
	}

	var dbOrder model.Order
	if err := sqlx.GetContext(ctx, ext, &dbOrder, ext.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotUpdated
		}

		return nil, fmt.Errorf("get context: %w", err)
	}

	return &dbOrder, nil
}

func insertOrderTimeline(ctx context.Context, ext sqlx.ExtContext, timeline model.OrderTimeline) error {
	if _, err := sqlx.NamedExecContext(ctx, ext, `
		INSERT INTO order_status_timeline(
			order_id,
			status,
			updated_at
		) VALUES(
			:order_id,
			:status,
			:updated_at
		)
	`, timeline,
	); err != nil {
		return fmt.Errorf("named exec context: %w", err)
	}

	return nil
}
