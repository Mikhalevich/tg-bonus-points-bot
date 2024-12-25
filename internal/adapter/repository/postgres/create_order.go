package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (p *Postgres) CreateOrder(ctx context.Context, coi port.CreateOrderInput) (order.ID, error) {
	query, args, err := sqlx.Named(`
		INSERT INTO orders(
			chat_id,
			status,
			verification_code,
			created_at
		) VALUES(
		 	:chat_id,
			:status,
			:verification_code,
			:created_at
		)
		RETURNING id
	`, map[string]any{
		"chat_id":           coi.ChatID.Int64(),
		"status":            coi.Status.String(),
		"verification_code": coi.VerificationCode,
		"created_at":        coi.StatusOperationTime,
	})

	if err != nil {
		return 0, fmt.Errorf("prepare named: %w", err)
	}

	var orderID int
	if err := sqlx.GetContext(ctx, p.db, &orderID, p.db.Rebind(query), args...); err != nil {
		if p.driver.IsConstraintError(err, "only_one_active_order_unique_idx") {
			return 0, errAlreadyExists
		}

		return 0, fmt.Errorf("get context: %w", err)
	}

	return order.IDFromInt(orderID), nil
}
