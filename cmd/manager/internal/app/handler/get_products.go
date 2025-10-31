package handler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/cmd/manager/internal/app/handler/request"
	"github.com/Mikhalevich/tg-bonus-points-bot/cmd/manager/internal/app/handler/response"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

type GetProductsInput struct {
	Filter request.Filter `query:"filter" enum:"enabled,disabled,all" required:"false" example:"enabled"`
}

type GetProductsOutput struct {
	Body response.Products
}

func (h *Handler) GetProducts(ctx context.Context, input *GetProductsInput) (*GetProductsOutput, error) {
	products, err := h.manager.GetProducts(ctx, product.Filter{
		Products: input.Filter.ToPortState(),
		Category: input.Filter.ToPortState(),
	})
	if err != nil {
		return nil, fmt.Errorf("get products: %w", err)
	}

	return &GetProductsOutput{
		Body: response.Products{
			Category: response.ConvertFromPortProduct(products),
		},
	}, nil
}
