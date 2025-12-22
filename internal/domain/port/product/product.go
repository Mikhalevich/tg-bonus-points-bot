package product

import (
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/flag"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/internal/id"
)

type ProductID struct {
	id.IntID
}

type CategoryID struct {
	id.IntID
}

func ProductIDFromInt(i int) ProductID {
	return ProductID{
		IntID: id.IntIDFromInt(i),
	}
}

func ProductIDFromString(s string) (ProductID, error) {
	intValue, err := id.IntIDFromString(s)
	if err != nil {
		return ProductID{}, fmt.Errorf("int id: %w", err)
	}

	return ProductID{
		IntID: intValue,
	}, nil
}

func CategoryIDFromInt(i int) CategoryID {
	return CategoryID{
		IntID: id.IntIDFromInt(i),
	}
}

func CategoryIDFromString(s string) (CategoryID, error) {
	intValue, err := id.IntIDFromString(s)
	if err != nil {
		return CategoryID{}, fmt.Errorf("int id: %w", err)
	}

	return CategoryID{
		IntID: intValue,
	}, nil
}

type Product struct {
	ID         ProductID
	Title      string
	CurrencyID currency.ID
	Price      int
	IsEnabled  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Category struct {
	ID        CategoryID
	Title     string
	IsEnabled bool
}

type CategoryProducts struct {
	ID       CategoryID
	Title    string
	Products []Product
}

type Filter struct {
	Products flag.State
	Category flag.State
}
