package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/store"
)

func (p *Postgres) GetStoreByID(ctx context.Context, id store.ID) (*store.Store, error) {
	modelStore, err := selectStoreByID(ctx, p.db, id)
	if err != nil {
		return nil, fmt.Errorf("select store by id: %w", err)
	}

	modelSchedule, err := selectStoreSchedule(ctx, p.db, id)
	if err != nil {
		return nil, fmt.Errorf("select store schedule: %w", err)
	}

	portStore, err := modelStore.ToPortStore(modelSchedule)
	if err != nil {
		return nil, fmt.Errorf("convert to port store: %w", err)
	}

	return portStore, nil
}

func selectStoreByID(ctx context.Context, tx sqlx.ExtContext, id store.ID) (*model.Store, error) {
	query, args, err := sqlx.Named(`
		SELECT
			id,
			description,
			default_currency_id
		FROM
			store
		WHERE
			id = :id
	`, map[string]any{
		"id": id.Int(),
	})

	if err != nil {
		return nil, fmt.Errorf("named: %w", err)
	}

	var s model.Store
	if err := sqlx.GetContext(ctx, tx, &s, tx.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}

		return nil, fmt.Errorf("get context: %w", err)
	}

	return &s, nil
}

func selectStoreSchedule(ctx context.Context, tx sqlx.ExtContext, id store.ID) ([]model.StoreSchedule, error) {
	query, args, err := sqlx.Named(`
		SELECT
			store_id,
			day_of_week,
			start_time,
			end_time
		FROM
			store_schedule
		WHERE
			store_id = :store_id
	`, map[string]any{
		"store_id": id.Int(),
	})

	if err != nil {
		return nil, fmt.Errorf("named: %w", err)
	}

	var schedule []model.StoreSchedule
	if err := sqlx.SelectContext(ctx, tx, &schedule, tx.Rebind(query), args...); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	return schedule, nil
}
