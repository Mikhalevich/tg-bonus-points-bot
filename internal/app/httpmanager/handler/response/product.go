package response

import (
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

type Products struct {
	Category []Category `json:"category"`
}

type Product struct {
	ID        int       `json:"id"`
	Title     string    `json:"titile"`
	Price     int       `json:"price"`
	IsEnabled bool      `json:"is_enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Category struct {
	Title    string    `json:"titile"`
	Products []Product `json:"products"`
}

func ConvertFromPortProduct(portCategory []product.Category) []Category {
	category := make([]Category, 0, len(portCategory))

	for _, v := range portCategory {
		products := make([]Product, 0, len(v.Products))
		for _, v := range v.Products {
			products = append(products, Product{
				ID:        v.ID,
				Title:     v.Title,
				Price:     v.Price,
				IsEnabled: v.IsEnabled,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
			})
		}

		category = append(category, Category{
			Title:    v.Title,
			Products: products,
		})
	}

	return category
}
