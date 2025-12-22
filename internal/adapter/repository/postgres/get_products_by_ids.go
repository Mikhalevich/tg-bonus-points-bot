package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func (p *Postgres) GetProductsByIDs(
	ctx context.Context,
	ids []product.ProductID,
	currencyID currency.ID,
) (map[product.ProductID]product.Product, error) {
	query, args, err := sqlx.In(`
		SELECT
			p.id,
			p.title,
			pp.currency_id,
			pp.price,
			p.is_enabled,
			p.created_at,
			p.updated_at
		FROM
			product p INNER JOIN product_price pp ON p.id = pp.product_id
		WHERE
			p.id IN(?) AND
			pp.currency_id = ?
	`, ids, currencyID.Int())

	if err != nil {
		return nil, fmt.Errorf("sqlx in: %w", err)
	}

	var dbProducts []model.Product
	if err := sqlx.SelectContext(ctx, p.db, &dbProducts, p.db.Rebind(query), args...); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	output := make(map[product.ProductID]product.Product, len(dbProducts))

	for _, v := range dbProducts {
		output[product.ProductIDFromInt(v.ID)] = v.ToPortProduct()
	}

	return output, nil
}
