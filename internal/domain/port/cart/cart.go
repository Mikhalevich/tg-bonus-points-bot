package cart

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

type ID string

func (id ID) String() string {
	return string(id)
}

func IDFromString(id string) ID {
	return ID(id)
}

type CartProduct struct {
	ProductID  product.ProductID
	CategoryID product.CategoryID
	Count      int
}
