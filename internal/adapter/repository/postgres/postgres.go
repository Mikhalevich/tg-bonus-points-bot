package postgres

import (
	"database/sql"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/jmoiron/sqlx"
)

var _ port.OrderRepository = (*Postgres)(nil)

type Postgres struct {
	db sqlx.ExtContext
}

func New(db *sql.DB, driverName string) *Postgres {
	return &Postgres{
		db: sqlx.NewDb(db, driverName),
	}
}
