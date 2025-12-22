package orderhistoryoffset

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	history "github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/orderhistory/v2"
)

var (
	_ history.Repository = (*OrderHistoryOffset)(nil)
)

type Driver interface {
	Name() string
	IsConstraintError(err error, constraint string) bool
}

type OrderHistoryOffset struct {
	db     sqlx.ExtContext
	driver Driver
}

func New(db *sql.DB, driver Driver) *OrderHistoryOffset {
	return &OrderHistoryOffset{
		db:     sqlx.NewDb(db, driver.Name()),
		driver: driver,
	}
}
