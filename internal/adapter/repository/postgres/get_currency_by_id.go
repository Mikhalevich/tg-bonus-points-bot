package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/currency"
)

func (p *Postgres) GetCurrencyByID(ctx context.Context, id currency.ID) (*currency.Currency, error) {
	curr, err := selectCurrencyByID(ctx, p.db, id)
	if err != nil {
		return nil, fmt.Errorf("select currency by id: %w", err)
	}

	return curr.ToPortCurrency(), nil
}

func selectCurrencyByID(ctx context.Context, trx sqlx.ExtContext, currencyID currency.ID) (*model.Currency, error) {
	query, args, err := sqlx.Named(`
		SELECT
			id,
			code,
			exp,
			decimal_sep,
			min_amount,
			max_amount,
			is_enabled
		FROM
			currency
		WHERE
			id = :id
	`, map[string]any{
		"id": currencyID.Int(),
	})

	if err != nil {
		return nil, fmt.Errorf("sqlx named: %w", err)
	}

	var curr model.Currency
	if err := sqlx.GetContext(ctx, trx, &curr, trx.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}

		return nil, fmt.Errorf("get context: %w", err)
	}

	return &curr, nil
}
