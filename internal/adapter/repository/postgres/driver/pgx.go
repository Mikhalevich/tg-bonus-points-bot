package driver

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/orderhistoryid"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/adapter/repository/postgres/orderhistoryoffset"
)

var (
	_ postgres.Driver           = (*Pgx)(nil)
	_ orderhistoryid.Driver     = (*Pgx)(nil)
	_ orderhistoryoffset.Driver = (*Pgx)(nil)
)

type Pgx struct {
}

func NewPgx() *Pgx {
	return &Pgx{}
}

func (p *Pgx) Name() string {
	return "pgx"
}

func (p *Pgx) IsConstraintError(err error, constraint string) bool {
	var pgxError *pgconn.PgError
	if !errors.As(err, &pgxError) {
		return false
	}

	if pgxError.ConstraintName == constraint {
		return true
	}

	return false
}
