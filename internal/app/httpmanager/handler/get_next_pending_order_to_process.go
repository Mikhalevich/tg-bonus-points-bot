package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/httpmanager/handler/response"
)

func (h *Handler) GetNextPendingOrderToProcess(w http.ResponseWriter, r *http.Request) error {
	inProgressOrder, err := h.manager.GetNextPendingOrderToProcess(r.Context())
	if err != nil {
		return fmt.Errorf("GetNextPendingOrderToProcess: %w", err)
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response.ToOrder(inProgressOrder)); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}

	return nil
}
