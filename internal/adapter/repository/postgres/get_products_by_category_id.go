package postgres

import (
	"context"
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

	return model.ToPortProducts(products), nil
}
