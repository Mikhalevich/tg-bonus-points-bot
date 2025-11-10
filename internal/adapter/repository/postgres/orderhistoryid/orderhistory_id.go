package orderhistoryid

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/customer/orderhistory"
)

var (
	_ orderhistory.Repository = (*OrderHistoryID)(nil)
)

type Driver interface {
	Name() string
	IsConstraintError(err error, constraint string) bool
}

type OrderHistoryID struct {
	db     sqlx.ExtContext
	driver Driver
}

func New(db *sql.DB, driver Driver) *OrderHistoryID {
	return &OrderHistoryID{
		db:     sqlx.NewDb(db, driver.Name()),
		driver: driver,
	}
}
