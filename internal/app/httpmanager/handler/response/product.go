package response

import (
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

type Products struct {
	Category []CategoryProducts `json:"category"`
}

type Product struct {
	ID        int       `json:"id"`
	Title     string    `json:"titile"`
	Price     int       `json:"price"`
	IsEnabled bool      `json:"is_enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryProducts struct {
	Title    string    `json:"titile"`
	Products []Product `json:"products"`
}

func ConvertFromPortProduct(portCategory []product.CategoryProducts) []CategoryProducts {
	categoryProducts := make([]CategoryProducts, 0, len(portCategory))

	for _, category := range portCategory {
		products := make([]Product, 0, len(category.Products))
		for _, prod := range category.Products {
			products = append(products, Product{
				ID:        prod.ID.Int(),
				Title:     prod.Title,
				Price:     prod.Price,
				IsEnabled: prod.IsEnabled,
				CreatedAt: prod.CreatedAt,
				UpdatedAt: prod.UpdatedAt,
			})
		}

		categoryProducts = append(categoryProducts, CategoryProducts{
			Title:    category.Title,
			Products: products,
		})
	}

	return categoryProducts
}
