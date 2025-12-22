package app

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/manager/internal/app/handler"
)

func (application *App) routes(handler *handler.Handler) {
	addRoute(application,
		huma.Operation{
			OperationID:   "next-order",
			Method:        http.MethodGet,
			Path:          "/order/next/",
			Summary:       "Get next order",
			Description:   "Get next order to process",
			Tags:          []string{"Order"},
			DefaultStatus: http.StatusOK,
		},
		handler.GetNextOrder,
	)

	addRoute(application,
		huma.Operation{
			OperationID:   "update-order-status",
			Method:        http.MethodPatch,
			Path:          "/order/{id}/",
			Summary:       "Update order status",
			Description:   "Update order status",
			Tags:          []string{"Order"},
			DefaultStatus: http.StatusOK,
		},
		handler.UpdateOrderStatus,
	)

	addRoute(application,
		huma.Operation{
			OperationID:   "get-products",
			Method:        http.MethodGet,
			Path:          "/product/",
			Summary:       "Get products",
			Description:   "Get available products",
			Tags:          []string{"Product"},
			DefaultStatus: http.StatusOK,
		},
		handler.GetProducts,
	)
}
