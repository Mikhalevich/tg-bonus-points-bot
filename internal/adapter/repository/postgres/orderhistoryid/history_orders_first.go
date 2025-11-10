package orderhistoryid

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (o *OrderHistoryID) HistoryOrdersFirst(
	ctx context.Context,
	chatID msginfo.ChatID,
	size int,
) ([]order.HistoryOrder, error) {
	query, args, err := sqlx.Named(`
		SELECT
			id,
			ROW_NUMBER() OVER (ORDER BY id) AS serial_number,
			status,
			currency_id,
			total_price,
			created_at
		FROM
			orders
		WHERE
			chat_id = :chat_id
		ORDER BY
			id DESC
		LIMIT
			:size
	`, map[string]any{
		"chat_id": chatID,
		"size":    size,
	})

	if err != nil {
		return nil, fmt.Errorf("sqlx named: %w", err)
	}

	orders, err := o.historyQuery(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("history query: %w", err)
	}

	return orders, nil
}

func (o *OrderHistoryID) historyQuery(ctx context.Context, query string, args ...any) ([]order.HistoryOrder, error) {
	var orders []model.HistoryOrder
	if err := sqlx.SelectContext(ctx, o.db, &orders, o.db.Rebind(query), args...); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	portHistoryOrders, err := model.ToPortHistoryOrders(orders)
	if err != nil {
		return nil, fmt.Errorf("convert to port orders: %w", err)
	}

	return portHistoryOrders, nil
}
