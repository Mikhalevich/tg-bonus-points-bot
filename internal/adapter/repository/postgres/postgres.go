package postgres

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/orderhistory"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port"
)

var (
	_ port.CustomerCartRepository         = (*Postgres)(nil)
	_ port.CustomerOrderPaymentRepository = (*Postgres)(nil)
	_ port.CustomerOrderActionRepository  = (*Postgres)(nil)
	_ port.ManagerRepository              = (*Postgres)(nil)
	_ port.StoreInfo                      = (*Postgres)(nil)

	_ orderhistory.CurrencyProvider = (*Postgres)(nil)
)

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
