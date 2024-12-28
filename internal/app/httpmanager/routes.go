package httpmanager

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/httpmanager/handler"
)

func (m *HTTPManager) routes(handler *handler.Handler) {
	m.addRoute("GET /order/next/", handler.GetNextPendingOrderToProcess)
	m.addRoute("PATCH /order/{id}/", handler.UpdateOrderStatus)
}
