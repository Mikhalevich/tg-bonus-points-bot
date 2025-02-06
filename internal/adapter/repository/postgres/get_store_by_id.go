package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/store"
	"github.com/jmoiron/sqlx"
)

func (p *Postgres) GetStoreByID(ctx context.Context, id store.ID) (*store.Store, error) {
	query, args, err := sqlx.Named(`
		SELECT
			id,
			description,
			default_currency_id,
		FROM
			store
		WHERE
			id = :id
	`, map[string]any{
		"id": id.Int(),
	})

	if err != nil {
		return nil, fmt.Errorf("named: %w", err)
	}

	var s model.Store
	if err := sqlx.GetContext(ctx, p.db, &s, p.db.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}

		return nil, fmt.Errorf("get context: %w", err)
	}

	return s.ToPortStore(), nil
}
