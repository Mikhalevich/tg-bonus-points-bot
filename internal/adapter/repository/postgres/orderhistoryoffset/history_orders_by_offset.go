package orderhistoryoffset

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (oh *OrderHistoryOffset) HistoryOrdersByOffset(
	ctx context.Context,
	chatID msginfo.ChatID,
	offset int,
	limit int,
) ([]order.HistoryOrder, error) {
	query, args, err := sqlx.Named(`
		SELECT
			id,
			ROW_NUMBER() OVER (ORDER BY id) AS serial_number,
			status,
			currency_id,
			created_at,
			total_price
		FROM
			orders
		WHERE
			chat_id = :chat_id
		ORDER BY
			id DESC
		LIMIT
			:limit
		OFFSET
			:offset
	`,
		map[string]any{
			"chat_id": chatID.Int64(),
			"offset":  offset,
			"limit":   limit,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("prepare query: %w", err)
	}

	var modelOrders []model.HistoryOrder
	if err := sqlx.SelectContext(ctx, oh.db, &modelOrders, oh.db.Rebind(query), args...); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	portOrders, err := model.ToPortHistoryOrders(modelOrders)
	if err != nil {
		return nil, fmt.Errorf("convert to port ordres: %w", err)
	}

	return portOrders, nil
}
