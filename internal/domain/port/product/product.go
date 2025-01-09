package product

import (
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/flag"
)

type ID int

func (id ID) Int() int {
	return int(id)
}

func IDFromInt(id int) ID {
	return ID(id)
}

type Product struct {
	ID        ID
	Title     string
	Price     int
	IsEnabled bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Category struct {
	ID       ID
	Title    string
	Products []Product
}

type Filter struct {
	Products flag.State
	Category flag.State
}
