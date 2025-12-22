package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/flag"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func (p *Postgres) GetCategoryProducts(ctx context.Context, filter product.Filter) ([]product.CategoryProducts, error) {
	var products []model.ProductCategory
	if err := sqlx.SelectContext(ctx, p.db, &products, buildProductsQuery(filter)); err != nil {
		return nil, fmt.Errorf("select products: %w", err)
	}

	return model.ToPortCategoryProducts(products), nil
}

func buildProductsQuery(filter product.Filter) string {
	return fmt.Sprintf(`
		SELECT
			p.id AS product_id,
			c.id AS category_id,
			p.title AS product_title,
			c.title AS category_title,
			p.price,
			p.is_enabled,
			p.created_at,
			p.updated_at
		FROM
			product_category pc INNER JOIN category c ON pc.category_id = c.id
			INNER JOIN product p ON pc.product_id = p.id
		%s
	`, buildProductsFilter(filter))
}

func buildProductsFilter(filter product.Filter) string {
	var conditions []string
	if filter.Products == flag.Enabled {
		conditions = append(conditions, "p.is_enabled = TRUE")
	}

	if filter.Products == flag.Disabled {
		conditions = append(conditions, "p.is_enabled = FALSE")
	}

	if filter.Category == flag.Enabled {
		conditions = append(conditions, "c.is_enabled = TRUE")
	}

	if filter.Category == flag.Disabled {
		conditions = append(conditions, "c.is_enabled = FALSE")
	}

	if len(conditions) == 0 {
		return ""
	}

	return fmt.Sprintf("WHERE %s", strings.Join(conditions, " AND "))
}
