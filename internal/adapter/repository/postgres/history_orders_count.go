package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

func (p *Postgres) HistoryOrdersCount(
	ctx context.Context,
	chatID msginfo.ChatID,
) (int, error) {
	query, args, err := sqlx.Named(`
		SELECT
			COUNT(*)
		FROM
			orders
		WHERE
			chat_id = :chat_id
	`, map[string]any{
		"chat_id": chatID,
	})

	if err != nil {
		return 0, fmt.Errorf("sqlx named: %w", err)
	}

	var count int
	if err := sqlx.GetContext(ctx, p.db, &count, p.db.Rebind(query), args...); err != nil {
		return 0, fmt.Errorf("select context: %w", err)
	}

	return count, nil
}
