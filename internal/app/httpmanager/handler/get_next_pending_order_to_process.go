package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/httpmanager/handler/response"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/logger"
)

func (h *Handler) GetNextPendingOrderToProcess(w http.ResponseWriter, r *http.Request) {
	inProgressOrder, err := h.manager.GetNextPendingOrderToProcess(r.Context())
	if err != nil {
		http.Error(w, err.Error(), convertErrorToHTTPCode(err))
		return
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response.ToOrder(inProgressOrder)); err != nil {
		logger.FromContext(r.Context()).
			WithError(err).
			Error("encode json response")
	}
}
