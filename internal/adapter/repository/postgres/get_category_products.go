package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (p *Postgres) GetCategoryProducts(ctx context.Context) ([]product.Category, error) {
	var products []model.ProductCategory
	if err := sqlx.SelectContext(ctx, p.db, &products, `
		SELECT
			p.id,
			p.title AS product_title,
			c.title AS category_title,
			p.price,
			p.is_enabled,
			p.created_at,
			p.updated_at
		FROM
			product_category pc INNER JOIN category c ON pc.category_id = c.id
			INNER JOIN product p ON pc.product_id = p.id
		WHERE
			p.is_enabled = TRUE AND
			c.is_enabled = TRUE
	`); err != nil {
		return nil, fmt.Errorf("select products: %w", err)
	}

	return model.ToPortCategory(products), nil
}
