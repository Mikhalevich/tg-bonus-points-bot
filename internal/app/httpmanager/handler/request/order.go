package request

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/httpmanager/internal/httperror"
)

type OrderStatus struct {
	Status string `json:"status"`
}

func (os OrderStatus) Validate() *httperror.ErrorHTTPResponse {
	if os.Status == "" {
		return httperror.BadRequest("status is the required field")
	}

	return nil
}
