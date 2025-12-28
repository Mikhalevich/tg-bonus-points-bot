package transaction

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type transactionCtxKey struct{}

type Transaction struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Transaction {
	return &Transaction{
		db: db,
	}
}

type TransactionFn func(ctx context.Context) error

func (t *Transaction) Transaction(ctx context.Context, trxFn TransactionFn) error {
	if activeTrx := trxFromContext(ctx); activeTrx != nil {
		if err := trxFn(ctx); err != nil {
			return fmt.Errorf("trx fn with active transaction: %w", err)
		}

		return nil
	}

	trx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin txx: %w", err)
	}

	//nolint:errcheck
	defer trx.Rollback()

	if err := trxFn(context.WithValue(ctx, transactionCtxKey{}, trx)); err != nil {
		return fmt.Errorf("trx fn: %w", err)
	}

	if err := trx.Commit(); err != nil {
		return fmt.Errorf("trx commit: %w", err)
	}

	return nil
}

func (t *Transaction) ExtContext(ctx context.Context) sqlx.ExtContext {
	if activeTrx := trxFromContext(ctx); activeTrx != nil {
		return activeTrx
	}

	return t.db
}

func trxFromContext(ctx context.Context) *sqlx.Tx {
	ctxValue := ctx.Value(transactionCtxKey{})

	if ctxValue == nil {
		return nil
	}

	activeTrx, ok := ctxValue.(*sqlx.Tx)
	if !ok {
		return nil
	}

	return activeTrx
}
