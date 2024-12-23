package postgres

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/jmoiron/sqlx"
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
		return 0, fmt.Errorf("get context: %w", err)
	}

	return order.IDFromInt(orderID), nil
}
