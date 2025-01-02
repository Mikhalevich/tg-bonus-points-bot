package product

import (
	"time"
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
