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

func (p Product) ToPortProduct() product.Product {
	return product.Product{
		ID:        product.IDFromInt(p.ID),
		Title:     p.Title,
		Price:     p.Price,
		IsEnabled: p.IsEnabled,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func ToPortProducts(dbProducts []Product) []product.Product {
	portProducts := make([]product.Product, 0, len(dbProducts))

	for _, p := range dbProducts {
		portProducts = append(portProducts, p.ToPortProduct())
	}

	return portProducts
}

type Category struct {
	ID        int    `db:"id"`
	Title     string `db:"title"`
	IsEnabled bool   `db:"is_enabled"`
}

type ProductCategory struct {
	ProductID     int       `db:"product_id"`
	CategoryID    int       `db:"category_id"`
	ProductTitle  string    `db:"product_title"`
	CategoryTitle string    `db:"category_title"`
	Price         int       `db:"price"`
	IsEnabled     bool      `db:"is_enabled"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (p ProductCategory) ToPortProduct() product.Product {
	return product.Product{
		ID:        product.IDFromInt(p.ProductID),
		Title:     p.ProductTitle,
		Price:     p.Price,
		IsEnabled: p.IsEnabled,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func ToPortCategory(dbProducts []ProductCategory) []product.Category {
	categoryMap := make(map[int]product.Category)
	for _, dbProduct := range dbProducts {
		if portProducts, ok := categoryMap[dbProduct.CategoryID]; ok {
			portProducts.Products = append(portProducts.Products, dbProduct.ToPortProduct())
			categoryMap[dbProduct.CategoryID] = portProducts

			continue
		}

		categoryMap[dbProduct.CategoryID] = product.Category{
			ID:       product.IDFromInt(dbProduct.CategoryID),
			Title:    dbProduct.CategoryTitle,
			Products: []product.Product{dbProduct.ToPortProduct()},
		}
	}

	category := make([]product.Category, 0, len(categoryMap))
	for _, v := range categoryMap {
		category = append(category, v)
	}

	return category
}
