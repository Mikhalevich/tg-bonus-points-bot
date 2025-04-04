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
	txFn TxFunc,
) error {
	trx, err := beginTx(ctx, s, allowNestedTransaction)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer trx.DeferCleanup()

	if err := txFn(ctx, trx); err != nil {
		if rollbackErr := trx.Rollback(); rollbackErr != nil {
			return errors.Join(fmt.Errorf("tx body: %w", err), fmt.Errorf("rollback: %w", err))
		}

		return fmt.Errorf("tx body: %w", err)
	}

	if err := trx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func beginTx(ctx context.Context, s sqlx.ExtContext, allowNestedTransaction bool) (*Tx, error) {
	dbConn, ok := s.(*sqlx.DB)
	if !ok {
		if allowNestedTransaction {
			if tx, ok := s.(*Tx); ok {
				return NewNestedTx(tx.Tx), nil
			}
		}

		return nil, errors.New("not sqlx db object")
	}

	tx, err := dbConn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	return NewTx(tx), nil
}
