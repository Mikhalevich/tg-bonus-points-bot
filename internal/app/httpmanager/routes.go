package httpmanager

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/httpmanager/handler"
)

func (m *HTTPManager) routes(handler *handler.Handler) {
	addRoute(m,
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

	addRoute(m,
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
}
