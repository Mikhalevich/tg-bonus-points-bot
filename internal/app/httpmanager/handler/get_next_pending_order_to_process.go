package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/httpmanager/handler/response"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/httpmanager/internal/httperror"
)

func (h *Handler) GetNextPendingOrderToProcess(w http.ResponseWriter, r *http.Request) *httperror.ErrorHTTPResponse {
	inProgressOrder, err := h.manager.GetNextPendingOrderToProcess(r.Context())
	if err != nil {
		return httperror.InternalServerError("GetNextPendingOrderToProcess", err)
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response.ToOrder(inProgressOrder)); err != nil {
		return httperror.InternalServerError("json encode", err)
	}

	return nil
}
