package product

import (
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/flag"
)

type Product struct {
	ID        int
	Title     string
	Price     int
	IsEnabled bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Category struct {
	Title    string
	Products []Product
}

type Filter struct {
	Products flag.State
	Category flag.State
}
