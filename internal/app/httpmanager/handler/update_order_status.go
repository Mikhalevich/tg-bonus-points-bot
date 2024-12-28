package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/httpmanager/handler/request"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/httpmanager/internal/httperror"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/perror"
)

func (h *Handler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) *httperror.ErrorHTTPResponse {
	var orderStatusReq request.OrderStatus
	if err := json.NewDecoder(r.Body).Decode(&orderStatusReq); err != nil {
		return httperror.InternalServerError("json decode request", err)
	}

	if herr := orderStatusReq.Validate(); herr != nil {
		return herr
	}

	orderID, err := order.IDFromString(r.PathValue("id"))
	if err != nil {
		return httperror.BadRequest("invalid id format").WithError(err)
	}

	status, err := order.StatusFromString(orderStatusReq.Status)
	if err != nil {
		return httperror.BadRequest("invalid status").WithError(err)
	}

	if err := h.manager.UpdateOrderStatus(r.Context(), orderID, status); err != nil {
		if perror.IsType(err, perror.TypeNotFound) {
			return httperror.NotFound("no orders with relevant status")
		}

		return httperror.InternalServerError("update order status", err)
	}

	return nil
}
