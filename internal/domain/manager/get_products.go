package manager

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (m *Manager) GetProducts(ctx context.Context, filter product.Filter) ([]product.CategoryProducts, error) {
	products, err := m.repository.GetCategoryProducts(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get category products: %w", err)
	}

	return products, nil
}
