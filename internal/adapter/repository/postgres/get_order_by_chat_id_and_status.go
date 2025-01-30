package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (p *Postgres) GetOrderByChatIDAndStatus(
	ctx context.Context,
	id msginfo.ChatID,
	statuses ...order.Status,
) (*order.Order, error) {
	dbOrder, err := selectOrderByChatIDAndStatus(ctx, p.db, id, statuses...)
	if err != nil {
		return nil, fmt.Errorf("select order by chat id and status: %w", err)
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

func selectOrderByChatIDAndStatus(
	ctx context.Context,
	ext sqlx.ExtContext,
	id msginfo.ChatID,
	statuses ...order.Status,
) (*model.Order, error) {
	query, args, err := sqlx.In(`
		SELECT
			id,
			chat_id,
			status,
			verification_code
		FROM
			orders
		WHERE
			chat_id = ? AND
			status IN(?)
	`, id, statuses)

	if err != nil {
		return nil, fmt.Errorf("sqlx.in: %w", err)
	}

	var dbOrder model.Order
	if err := sqlx.GetContext(ctx, ext, &dbOrder, ext.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}

		return nil, fmt.Errorf("get context: %w", err)
	}

	return &dbOrder, nil
}

func selectOrderTimeline(ctx context.Context, ext sqlx.ExtContext, orderID int) ([]model.OrderTimeline, error) {
	query, args, err := sqlx.Named(`
		SELECT
			order_id,
			status,
			updated_at
		FROM
			order_status_timeline
		WHERE
			order_id = :order_id
	`, map[string]any{
		"order_id": orderID,
	})

	if err != nil {
		return nil, fmt.Errorf("named: %w", err)
	}

	var orderTimeline []model.OrderTimeline
	if err := sqlx.SelectContext(ctx, ext, &orderTimeline, ext.Rebind(query), args...); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	return orderTimeline, nil
}

func selectOrderProducts(ctx context.Context, ext sqlx.ExtContext, id int) ([]model.OrderProduct, error) {
	query, args, err := sqlx.Named(`
		SELECT
			order_id,
			product_id,
			count,
			price
		FROM
			order_products
		WHERE
			order_id = :order_id
	`, map[string]any{
		"order_id": id,
	})

	if err != nil {
		return nil, fmt.Errorf("named: %w", err)
	}

	var orderProducts []model.OrderProduct
	if err := sqlx.SelectContext(ctx, ext, &orderProducts, ext.Rebind(query), args...); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	return orderProducts, nil
}
