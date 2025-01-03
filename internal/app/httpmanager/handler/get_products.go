package handler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/httpmanager/handler/response"
)

type GetProductsInput struct {
}

type GetProductsOutput struct {
	Body response.Products
}

func (h *Handler) GetProducts(ctx context.Context, input *GetProductsInput) (*GetProductsOutput, error) {
	products, err := h.manager.GetProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("get products: %w", err)
	}

	return &GetProductsOutput{
		Body: response.Products{
			Category: response.ConvertFromPortProduct(products),
		},
	}, nil
}
