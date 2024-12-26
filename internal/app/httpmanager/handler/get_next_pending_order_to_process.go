package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) GetNextPendingOrderToProcess(w http.ResponseWriter, r *http.Request) {
	_, err := h.manager.GetNextPendingOrderToProcess(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "it works...")
}
