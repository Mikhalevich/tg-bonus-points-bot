package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/transaction"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (p *Postgres) CreateOrder(ctx context.Context, coi port.CreateOrderInput) (order.ID, error) {
	var (
		orderID order.ID
		err     error
	)

	if err := transaction.Transaction(ctx, p.db, true, func(ctx context.Context, tx sqlx.ExtContext) error {
		orderID, err = p.insertOrder(ctx, tx, model.Order{
			ChatID:           coi.ChatID.Int64(),
			Status:           coi.Status.String(),
			VerificationCode: coi.VerificationCode,
		})

		if err != nil {
			return fmt.Errorf("insert order: %w", err)
		}

		if err := insertOrderTimeline(ctx, tx, model.OrderTimeline{
			ID:        orderID.Int(),
			Status:    coi.Status.String(),
			UpdatedAt: coi.StatusOperationTime,
		}); err != nil {
			return fmt.Errorf("insert order timeline: %w", err)
		}

		return nil
	}); err != nil {
		return 0, fmt.Errorf("transaction: %w", err)
	}

	return orderID, nil
}

func (p *Postgres) insertOrder(ctx context.Context, ext sqlx.ExtContext, dbOrder model.Order) (order.ID, error) {
	query, args, err := sqlx.Named(`
		INSERT INTO orders(
			chat_id,
			status,
			verification_code
		) VALUES (
			:chat_id,
			:status,
			:verification_code
		)
		RETURNING id
	`, dbOrder)

	if err != nil {
		return 0, fmt.Errorf("prepare named: %w", err)
	}

	var orderID int
	if err := sqlx.GetContext(ctx, ext, &orderID, ext.Rebind(query), args...); err != nil {
		if p.driver.IsConstraintError(err, "only_one_active_order_unique_idx") {
			return 0, errAlreadyExists
		}

		return 0, fmt.Errorf("insert order: %w", err)
	}

	return order.IDFromInt(orderID), nil
}
