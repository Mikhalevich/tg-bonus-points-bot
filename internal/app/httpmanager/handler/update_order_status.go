package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/httpmanager/handler/request"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (h *Handler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) error {
	var orderStatusReq request.OrderStatus
	if err := json.NewDecoder(r.Body).Decode(&orderStatusReq); err != nil {
		return fmt.Errorf("json decode request: %w", err)
	}

	if err := orderStatusReq.Validate(); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	orderID, err := order.IDFromString(orderStatusReq.ID)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	status, err := order.StatusFromString(orderStatusReq.Status)
	if err != nil {
		return fmt.Errorf("invalid status: %w", err)
	}

	if err := h.manager.UpdateOrderStatus(r.Context(), orderID, status); err != nil {
		return fmt.Errorf("UpdateOrderStatus: %w", err)
	}

	return nil
}
