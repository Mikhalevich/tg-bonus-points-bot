package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (p *Postgres) GetProductsByCategoryID(ctx context.Context, id product.ID) ([]product.Product, error) {
	query, args, err := sqlx.Named(`
		SELECT
			p.id,
			p.title,
			p.price,
			p.is_enabled,
			p.created_at,
			p.updated_at
		FROM
			product p INNER JOIN product_category pc ON p.id = pc.product_id
		WHERE
			pc.category_id = :category_id AND
			p.is_enabled = :is_product_enabled
		ORDER BY
			p.title
	`, map[string]any{
		"category_id":        id.Int(),
		"is_product_enabled": true,
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
