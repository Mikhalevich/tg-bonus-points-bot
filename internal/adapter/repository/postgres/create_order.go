package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/transaction"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (p *Postgres) CreateOrder(ctx context.Context, coi port.CreateOrderInput) (*order.Order, error) {
	var orderResult order.Order

	if err := transaction.Transaction(ctx, p.db, true, func(ctx context.Context, tx sqlx.ExtContext) error {
		orderID, err := p.insertOrder(ctx, tx, model.Order{
			ChatID:           coi.ChatID.Int64(),
			Status:           coi.Status.String(),
			VerificationCode: coi.VerificationCode,
			CurrencyID:       coi.CurrencyID.Int(),
		})

		if err != nil {
			return fmt.Errorf("insert order: %w", err)
		}

		if err := insertProductsToOrder(ctx, tx, model.PortToOrderProducts(orderID, coi.Products)); err != nil {
			return fmt.Errorf("insert order products: %w", err)
		}

		if err := insertOrderTimeline(ctx, tx, model.OrderTimeline{
			ID:        orderID.Int(),
			Status:    coi.Status.String(),
			UpdatedAt: coi.StatusOperationTime,
		}); err != nil {
			return fmt.Errorf("insert order timeline: %w", err)
		}

		orderResult = convertToOrder(orderID, coi)

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %w", err)
	}

	return &orderResult, nil
}

func convertToOrder(id order.ID, input port.CreateOrderInput) order.Order {
	return order.Order{
		ID:               id,
		ChatID:           input.ChatID,
		Status:           input.Status,
		VerificationCode: input.VerificationCode,
		CurrencyID:       input.CurrencyID,
		Products:         input.Products,
	}
}

func (p *Postgres) insertOrder(ctx context.Context, ext sqlx.ExtContext, dbOrder model.Order) (order.ID, error) {
	query, args, err := sqlx.Named(`
		INSERT INTO orders(
			chat_id,
			status,
			verification_code,
			currency_id
		) VALUES (
			:chat_id,
			:status,
			:verification_code,
			:currency_id
		)
		RETURNING id
	`, dbOrder)

	if err != nil {
		return 0, fmt.Errorf("prepare named: %w", err)
	}

	var orderID int
	if err := sqlx.GetContext(ctx, ext, &orderID, ext.Rebind(query), args...); err != nil {
		if p.driver.IsConstraintError(err, "orders_only_one_active_order_unique_idx") {
			return 0, errAlreadyExists
		}

		return 0, fmt.Errorf("insert order: %w", err)
	}

	return order.IDFromInt(orderID), nil
}

func insertProductsToOrder(ctx context.Context, ext sqlx.ExtContext, products []model.OrderProduct) error {
	res, err := sqlx.NamedExecContext(ctx, ext, `
		INSERT INTO order_products(
			order_id,
			product_id,
			count,
			price
		) VALUES (
			:order_id,
			:product_id,
			:count,
			:price
		)`,
		products,
	)

	if err != nil {
		return fmt.Errorf("insert order products: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return errNotUpdated
	}

	return nil
}
