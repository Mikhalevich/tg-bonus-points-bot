package postgres

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
)

var _ port.OrderRepository = (*Postgres)(nil)

type Driver interface {
	Name() string
	IsConstraintError(err error, constraint string) bool
}

type Postgres struct {
	db     sqlx.ExtContext
	driver Driver
}

func New(db *sql.DB, driver Driver) *Postgres {
	return &Postgres{
		db:     sqlx.NewDb(db, driver.Name()),
		driver: driver,
	}
}
