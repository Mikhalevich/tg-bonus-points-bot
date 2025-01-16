package cart

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

type ID string

func (id ID) String() string {
	return string(id)
}

func IDFromString(id string) ID {
	return ID(id)
}

type CartProduct struct {
	Product product.Product
	Count   int
}
