package transaction

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type SqlxDB struct {
	db *sqlx.DB
}

func NewSqlxDB(db *sqlx.DB) SqlxDB {
	return SqlxDB{
		db: db,
	}
}

func (s SqlxDB) Begin(ctx context.Context) (DBTx, error) {
	trx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin txx: %w", err)
	}

	return NewSqlxDBTx(trx), nil
}

type SqlxDBTx struct {
	*sqlx.Tx
}

func NewSqlxDBTx(trx *sqlx.Tx) SqlxDBTx {
	return SqlxDBTx{
		Tx: trx,
	}
}
