package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (p *Postgres) GetProductsByCategoryID(
	ctx context.Context,
	categoryID product.CategoryID,
	currencyID currency.ID,
) ([]product.Product, error) {
	cur, err := selectCurrencyByID(ctx, p.db, currencyID)
	if err != nil {
		return nil, fmt.Errorf("currency by id: %w", err)
	}

	query, args, err := sqlx.Named(`
		SELECT
			p.id,
			p.title,
			pp.currency_id,
			pp.price,
			p.is_enabled,
			p.created_at,
			p.updated_at
		FROM
			product p INNER JOIN product_category pc ON p.id = pc.product_id
			INNER JOIN product_price pp ON p.id = pp.product_id
		WHERE
			pc.category_id = :category_id AND
			p.is_enabled = :is_product_enabled AND
			pp.currency_id = :currency_id
		ORDER BY
			p.title
	`, map[string]any{
		"category_id":        categoryID.Int(),
		"is_product_enabled": true,
		"currency_id":        currencyID.Int(),
	})

	if err != nil {
		return nil, fmt.Errorf("prepare named: %w", err)
	}

	var products []model.Product
	if err := sqlx.SelectContext(ctx, p.db, &products, p.db.Rebind(query), args...); err != nil {
		return nil, fmt.Errorf("select products: %w", err)
	}

	return model.ToPortProducts(products, cur), nil
}

func selectCurrencyByID(ctx context.Context, tx sqlx.ExtContext, id currency.ID) (model.Currency, error) {
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
		"id": id.Int(),
	})

	if err != nil {
		return model.Currency{}, fmt.Errorf("sqlx named: %w", err)
	}

	var cur model.Currency
	if err := sqlx.GetContext(ctx, tx, &cur, tx.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Currency{}, errNotFound
		}

		return model.Currency{}, fmt.Errorf("get context: %w", err)
	}

	return cur, nil
}
