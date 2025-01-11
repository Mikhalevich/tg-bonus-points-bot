package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (p *Postgres) GetProductsByIDs(
	ctx context.Context,
	ids []product.ID,
) (map[product.ID]product.Product, error) {
	query, args, err := sqlx.In(`
		SELECT
			id,
			title,
			price,
			is_enabled,
			created_at,
			updated_at
		FROM
			product
		WHERE id IN(?)
	`, ids)

	if err != nil {
		return nil, fmt.Errorf("sqlx in: %w", err)
	}

	var dbProducts []model.Product
	if err := sqlx.SelectContext(ctx, p.db, &dbProducts, p.db.Rebind(query), args...); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	output := make(map[product.ID]product.Product, len(dbProducts))

	for _, v := range dbProducts {
		output[product.IDFromInt(v.ID)] = v.ToPortProduct()
	}

	return output, nil
}
