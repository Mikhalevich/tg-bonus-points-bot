package model

import (
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

type Product struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	Price     int       `db:"price"`
	IsEnabled bool      `db:"is_enabled"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Category struct {
	ID        int    `db:"id"`
	Title     string `db:"title"`
	IsEnabled bool   `db:"is_enabled"`
}

type ProductCategory struct {
	ID            int       `db:"id"`
	ProductTitle  string    `db:"product_title"`
	CategoryTitle string    `db:"category_title"`
	Price         int       `db:"price"`
	IsEnabled     bool      `db:"is_enabled"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (p ProductCategory) ToPortProduct() product.Product {
	return product.Product{
		ID:        p.ID,
		Title:     p.ProductTitle,
		Price:     p.Price,
		IsEnabled: p.IsEnabled,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func ToPortCategory(dbProducts []ProductCategory) []product.Category {
	categoryMap := make(map[string][]product.Product)
	for _, dbProduct := range dbProducts {
		if portProducts, ok := categoryMap[dbProduct.CategoryTitle]; ok {
			categoryMap[dbProduct.CategoryTitle] = append(portProducts, dbProduct.ToPortProduct())
			continue
		}

		categoryMap[dbProduct.CategoryTitle] = []product.Product{dbProduct.ToPortProduct()}
	}

	category := make([]product.Category, 0, len(categoryMap))
	for k, v := range categoryMap {
		category = append(category, product.Category{
			Title:    k,
			Products: v,
		})
	}

	return category
}
