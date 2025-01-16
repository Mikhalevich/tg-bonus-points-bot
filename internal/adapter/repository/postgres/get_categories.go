package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (p *Postgres) GetCategories(ctx context.Context) ([]product.Category, error) {
	var (
		query = `
			SELECT
				id,
				title,
				is_enabled
			FROM
				category
			WHERE
				is_enabled = TRUE AND
				id IN (
					SELECT
						pc.category_id
					FROM
						product_category pc INNER JOIN product p ON pc.product_id = p.id
					WHERE
						p.is_enabled = TRUE
					GROUP BY
						pc.category_id
					
				)
			ORDER BY
				title
			`
		categories []model.Category
	)

	if err := sqlx.SelectContext(ctx, p.db, &categories, query); err != nil {
		return nil, fmt.Errorf("select categories: %w", err)
	}

	return model.ToPortCategories(categories), nil
}
