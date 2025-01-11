package product

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/flag"
)

type ID int

func (id ID) Int() int {
	return int(id)
}

func (id ID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

func IDFromInt(id int) ID {
	return ID(id)
}

func IDFromString(id string) (ID, error) {
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse int: %w", err)
	}

	return ID(intID), nil
}

type Product struct {
	ID        ID
	Title     string
	Price     int
	IsEnabled bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProductCount struct {
	Product Product
	Count   int
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
