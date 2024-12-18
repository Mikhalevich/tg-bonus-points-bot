package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TxFunc func(ctx context.Context, tx sqlx.ExtContext) error

func Transaction(
	ctx context.Context,
	s sqlx.ExtContext,
	allowNestedTransaction bool,
	fn TxFunc,
) error {
	tx, err := beginTx(ctx, s, allowNestedTransaction)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		tx.DeferCleanup()
	}()

	if err := fn(ctx, tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(fmt.Errorf("tx body: %w", err), fmt.Errorf("rollback: %w", err))
		}

		return fmt.Errorf("tx body: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func beginTx(ctx context.Context, s sqlx.ExtContext, allowNestedTransaction bool) (*Tx, error) {
	db, ok := s.(*sqlx.DB)
	if !ok {
		if allowNestedTransaction {
			if tx, ok := s.(*Tx); ok {
				return NewNestedTx(tx.Tx), nil
			}
		}

		return nil, errors.New("not sqlx db object")
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	return NewTx(tx), nil
}
