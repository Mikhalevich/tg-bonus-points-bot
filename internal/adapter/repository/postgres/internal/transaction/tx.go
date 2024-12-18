package transaction

import (
	"github.com/jmoiron/sqlx"
)

type Tx struct {
	*sqlx.Tx
	IsNested bool
}

func NewTx(tx *sqlx.Tx) *Tx {
	return &Tx{
		Tx: tx,
	}
}

func NewNestedTx(tx *sqlx.Tx) *Tx {
	return &Tx{
		Tx:       tx,
		IsNested: true,
	}
}

func (t *Tx) Commit() error {
	if t.IsNested {
		return nil
	}

	//nolint:wrapcheck
	return t.Tx.Commit()
}

func (t *Tx) Rollback() error {
	//nolint:wrapcheck
	return t.Tx.Rollback()
}

func (t *Tx) DeferCleanup() {
	if t.IsNested {
		return
	}

	//nolint:errcheck
	t.Tx.Rollback()
}
