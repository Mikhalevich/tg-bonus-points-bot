package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (p *Postgres) GetOrderByID(ctx context.Context, id order.ID) (*order.Order, error) {
	dbOrder, err := selectOrderByID(ctx, p.db, id)
	if err != nil {
		return nil, fmt.Errorf("select order by id: %w", err)
	}

	orderProducts, err := selectOrderProducts(ctx, p.db, dbOrder.ID)
	if err != nil {
		return nil, fmt.Errorf("select order products: %w", err)
	}

	orderTimeline, err := selectOrderTimeline(ctx, p.db, dbOrder.ID)
	if err != nil {
		return nil, fmt.Errorf("select order timeline: %w", err)
	}

	portOrder, err := model.ToPortOrder(dbOrder, orderProducts, orderTimeline)
	if err != nil {
		return nil, fmt.Errorf("convert to port order: %w", err)
	}

	return portOrder, nil
}

func selectOrderByID(ctx context.Context, ext sqlx.ExtContext, id order.ID) (*model.Order, error) {
	var modelOrder model.Order
	if err := sqlx.GetContext(ctx, ext, &modelOrder,
		`SELECT
			id,
			chat_id,
			status,
			verification_code,
			currency_id
		FROM
			orders
		WHERE
			id = $1
	`, id.Int(),
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}

		return nil, fmt.Errorf("get context: %w", err)
	}

	return &modelOrder, nil
}
